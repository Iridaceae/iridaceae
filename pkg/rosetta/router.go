package rosetta

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/Iridaceae/iridaceae/pkg/helpers"

	"github.com/Iridaceae/iridaceae/pkg/log"

	"github.com/bwmarrin/discordgo"

	"github.com/sarulabs/di/v2"
)

// SpliceRegex represents the regex to split arguments.
var SpliceRegex = regexp.MustCompile(`\\s+`)

// StdRouter can be used out of the box, acts as a master router for iris.
var StdRouter Router

func init() {
	StdRouter = NewRouter(NewDefaultConfig())
}

// Config setup configs value for our router.
type Config struct {
	GeneralPrefix         string `json:"general_prefix"`
	IgnoreCase            bool   `json:"ignore_case"`
	AllowDM               bool   `json:"allow_dm"`
	AllowBots             bool   `json:"allow_bots"`
	ExecuteOnEdit         bool   `json:"execute_on_edit"`
	UseDefaultHelpCommand bool   `json:"user_default_help_command"`
	DeleteMessageAfter    bool   `json:"delete_message_after"`

	// ObjectContainer can be passed by user to obtain instances from context.
	ObjectContainer di.Container `json:"-"`

	// OnError will be called when router failed to execute the command.
	// OnError will be passed when context failed to run, and return an ErrorType and error objects.
	OnError func(ctx Context, errType ErrorType, err error)

	// GuildPrefixGetter is called to get guild prefix.
	// Function will have guild id passed and return
	// the guild prefix if specified, else it will return
	// default prefix.
	// An error will be returned when the retrieving of the guild prefix failed unexpectedly.
	GuildPrefixGetter func(gid string) (string, error)
}

// Router defines a command register and muxer.
type Router interface {
	ReadOnlyObjectMap

	// Register is shortened for RegisterMiddleware and RegisterCommand
	// and automatically chooses depending on implementation.
	//
	// panics if an instance is passed which neither implements Command and Middleware.
	Register(v interface{})

	// RegisterCommand registers the passed Command interface.
	RegisterCommand(cmd Command)

	// RegisterMiddleware registers Middleware interface.
	RegisterMiddleware(m Middleware)

	// Setup registers given handlers to the passed discordgo.Session which are
	// used to handle and parse command.
	Setup(s *discordgo.Session)

	// GetConfig returns the specified config object which was specified on initialization.
	GetConfig() *Config

	// GetCommandMap returns internal command map.
	GetCommandMap() map[string]Command

	// GetCommandInstances returns an array of all registered command instance.
	GetCommandInstances() []Command

	// GetCommand returns a command instance from the registry by invoker. If command could
	// not be found, false is returned.
	GetCommand(invoke string) (Command, bool)
}

// router is our default implementation of Router.

type router struct {
	config          *Config
	cmdMap          map[string]Command
	cmdInstances    []Command
	middleware      []Middleware
	objectContainer di.Container
	ctxPool         *sync.Pool
	objectMap       *sync.Map
}

func NewDefaultConfig() *Config {
	return &Config{
		GeneralPrefix:         "r!",
		IgnoreCase:            true,
		AllowDM:               true,
		AllowBots:             false,
		ExecuteOnEdit:         true,
		UseDefaultHelpCommand: false,
		DeleteMessageAfter:    true,
		OnError: func(ctx Context, errType ErrorType, err error) {
			msg, _ := ctx.RespondEmbedError(getErrorTypeName(errType), err)
			log.Error(err).Msgf("username: %s#%s sent:%s error: %+v", ctx.GetUser().Username, ctx.GetUser().Discriminator, ctx.GetMessage().Content, err)
			helpers.DeleteMessageAfter(ctx.GetSession(), msg, 60*time.Second)
		},
	}
}

func NewRouter(c *Config) Router {
	if c.OnError == nil {
		// setup a default onerror func.
		c.OnError = func(ctx Context, errType ErrorType, err error) {}
	}
	if c.GuildPrefixGetter == nil {
		c.GuildPrefixGetter = func(string) (string, error) { return "", nil }
	}
	r := &router{
		config:          c,
		cmdMap:          make(map[string]Command),
		cmdInstances:    make([]Command, 0),
		objectContainer: c.ObjectContainer,
		ctxPool:         &sync.Pool{New: func() interface{} { return &context{objectMap: &sync.Map{}} }},
		objectMap:       &sync.Map{},
	}

	if r.objectContainer == nil {
		builder, _ := di.NewBuilder()
		r.objectContainer = builder.Build()
	}

	if c.UseDefaultHelpCommand {
		r.RegisterCommand(&DefaultHelpCommand{})
	}
	return r
}

func (r *router) GetObject(key string) interface{} {
	value, err := r.objectContainer.SafeGet(key)
	if err != nil {
		value, _ = r.objectMap.Load(key)
	}
	return value
}

func (r *router) SetObject(key string, value interface{}) {
	r.objectMap.Store(key, value)
}

func (r *router) Register(v interface{}) {
	switch i := v.(type) {
	case Command:
		r.RegisterCommand(i)
	case Middleware:
		r.RegisterMiddleware(i)
	default:
		panic("instance doesn't implements Command or Middleware")
	}
}

func (r *router) RegisterCommand(cmd Command) {
	r.cmdInstances = append(r.cmdInstances, cmd)
	for _, i := range cmd.GetInvokers() {
		if r.config.IgnoreCase {
			i = strings.ToLower(i)
		}
		if _, ok := r.cmdMap[i]; ok {
			panic(fmt.Sprintf("invoke %s already registered, panicked!", i))
		}
		r.cmdMap[i] = cmd
	}
}

func (r *router) RegisterMiddleware(m Middleware) {
	r.middleware = append(r.middleware, m)
}

func (r *router) Setup(s *discordgo.Session) {
	s.AddHandler(func(s *discordgo.Session, e *discordgo.MessageCreate) {
		r.trigger(s, e.Message, false)
	})
	if r.config.ExecuteOnEdit {
		s.AddHandler(func(s *discordgo.Session, e *discordgo.MessageUpdate) {
			r.trigger(s, e.Message, false)
		})
	}
}

func (r *router) trigger(s *discordgo.Session, msg *discordgo.Message, edit bool) {
	var err error
	prefix := ""

	// check if given message author is a bot.
	if !r.config.AllowBots || msg.Author == nil || msg.Author.Bot || msg.Author.ID == s.State.User.ID {
		return
	}

	ctx, _ := r.ctxPool.Get().(*context)
	ctx.router = r
	ctx.session = s
	ctx.message = msg
	ctx.member = msg.Member
	ctx.isEdit = edit

	defer func() {
		clearMap(ctx.objectMap)
		r.ctxPool.Put(ctx)
	}()

	content, containsPrefix := hasPrefix(msg.Content, r.config.GeneralPrefix, r.config.IgnoreCase)
	if containsPrefix {
		prefix = r.config.GeneralPrefix
	} else {
		var guildPrefix string
		guildPrefix, err = r.config.GuildPrefixGetter(msg.GuildID)
		if err != nil {
			r.config.OnError(ctx, ErrTypeGuildPrefixGetter, err)
			return
		}
		_, ok := hasPrefix(msg.Content, guildPrefix, r.config.IgnoreCase)
		if guildPrefix != "" && ok {
			prefix = guildPrefix
		}
	}

	content = strings.TrimSpace(content)
	// if message is empty after prefix processing then we don't do anything.
	if content == "" || prefix == "" && ctx.channel.Type != discordgo.ChannelTypeDM {
		return
	}

	if ctx.channel, err = s.State.Channel(msg.ChannelID); err != nil {
		if ctx.channel, err = s.Channel(msg.ChannelID); err != nil {
			r.config.OnError(ctx, ErrTypeGetChannel, err)
			return
		}
	}

	ctx.isDM = ctx.channel.Type == discordgo.ChannelTypeDM || ctx.channel.Type == discordgo.ChannelTypeGroupDM
	if !r.config.AllowDM && ctx.isDM {
		return
	}

	if !ctx.isDM {
		if ctx.guild, err = s.State.Guild(msg.GuildID); err != nil {
			if ctx.guild, err = s.Guild(msg.GuildID); err != nil {
				r.config.OnError(ctx, ErrTypeGetGuild, err)
				return
			}
		}
	}

	args := ParseArguments(msg.Content)
	invoke, arg := args.Args()[0].String(), args.Args()[1:]
	ctx.args = FromArguments(arg)

	cmd, ok := r.GetCommand(invoke)
	if !ok {
		r.config.OnError(ctx, ErrTypeCommandNotFound, ErrCommandNotFound)
	}

	if ctx.isDM && !cmd.IsExecutableInDM() {
		r.config.OnError(ctx, ErrTypeNotExecutableInDM, ErrNotExecutableInDMs)
	}

	ctx.SetObject(ObjectMapKeyRouter, r)

	if !r.executeMiddlewares(cmd, ctx, LayerBeforeCommand) {
		return
	}

	if err = cmd.Exec(ctx); err != nil {
		r.config.OnError(ctx, ErrTypeCommandExec, err)
		return
	}

	if !r.executeMiddlewares(cmd, ctx, LayerAfterCommand) {
		return
	}

	if r.config.DeleteMessageAfter {
		if err = s.ChannelMessageDelete(msg.ChannelID, msg.ID); err != nil {
			r.config.OnError(ctx, ErrTypeDeleteCommandMessage, err)
			return
		}
	}
}

func (r *router) executeMiddlewares(cmd Command, ctx Context, layer MiddlewareLayer) bool {
	for _, m := range r.middleware {
		if m.GetLayer()&layer == 0 {
			continue
		}

		next, err := m.Handle(cmd, ctx, layer)
		if err != nil {
			r.config.OnError(ctx, ErrTypeMiddleware, err)
			return false
		}
		if !next {
			return false
		}
	}
	return true
}

func (r *router) GetConfig() *Config {
	return r.config
}

func (r *router) GetCommandMap() map[string]Command {
	return r.cmdMap
}

func (r *router) GetCommandInstances() []Command {
	return r.cmdInstances
}

func (r *router) GetCommand(invoke string) (Command, bool) {
	if r.config.IgnoreCase {
		invoke = strings.ToLower(invoke)
	}
	cmd, ok := r.cmdMap[invoke]
	return cmd, ok
}
