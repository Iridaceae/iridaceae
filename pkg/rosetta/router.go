package rosetta

import (
	"regexp"
	"strings"

	"github.com/Iridaceae/iridaceae/pkg/stlog"

	"github.com/bwmarrin/discordgo"
)

var (
	// SpliceRegex represents the regex to split arguments.
	SpliceRegex = regexp.MustCompile(`\\s+`)
	// StdRouter can be used out of the box, acts as a master router for iris.
	StdRouter *Router
)

func init() {
	StdRouter = Create(&Router{
		Prefixes:         []string{"!", "-", "ir-"},
		IgnorePrefixCase: true,
		BotsAllowed:      false,
		Logger:           stlog.Defaults,
		Commands:         []*Command{},
		Middlewares:      []Middleware{},
		PingHandler: func(ctx *Context) {
			// TODO: Default PingHandler should returns dog facts or a paragraph from GPT3 =))
			if err := ctx.RespondText("Hello World"); err != nil {
				panic(err)
			}
		},
	})
}

// Router represents a discordgo command routers.
// A derivation of Mux from disgord (https://github.com/bwmarrin/disgord).
type Router struct {
	Prefixes         []string
	IgnorePrefixCase bool
	BotsAllowed      bool
	Logger           *stlog.Logger
	Commands         []*Command
	Middlewares      []Middleware
	PingHandler      ExecutionHandler
	Storage          map[string]*ObjectsMap
}

// Create ensures that router storage map is initialized.
func Create(r *Router) *Router {
	r.Storage = make(map[string]*ObjectsMap)
	r.Logger = stlog.Defaults
	return r
}

// RegisterCmd adds a new commands to routers.
func (r *Router) RegisterCmd(cmd *Command) {
	r.Commands = append(r.Commands, cmd)
}

// GetCmd returns command with given name if exists.
func (r *Router) GetCmd(name string) *Command {
	for _, cmd := range r.Commands {
		toCheck := getIdentifiers(cmd)

		// check prefix of string.
		if arrayContains(toCheck, name, cmd.IgnoreCase) {
			return cmd
		}
	}
	return nil
}

// RegisterMiddleware registers a new middleware.
func (r *Router) RegisterMiddleware(middleware Middleware) {
	r.Middlewares = append(r.Middlewares, middleware)
}

// InitializeStorage initializes a storage map.
func (r *Router) InitializeStorage(name string) {
	r.Storage[name] = newObjectsMap()
}

// Initialize discordgo message even listener.
func (r *Router) Initialize(s *discordgo.Session) {
	s.AddHandler(r.Handler())
}

// Handler provides discordgo handler for given router.
func (r *Router) Handler() func(*discordgo.Session, *discordgo.MessageCreate) {
	return func(session *discordgo.Session, event *discordgo.MessageCreate) {
		message := event.Message
		content := message.Content

		// similar to onMessageReceived, check if message was sent by a bot.
		if message.Author.Bot && !r.BotsAllowed {
			return
		}

		// Execute ping handler if message equals the current bot's mentions.
		if (content == "<@!"+session.State.User.ID+">" || content == "<@"+session.State.User.ID+">") && r.PingHandler != nil {
			r.PingHandler(&Context{
				Session:   session,
				Event:     event,
				Arguments: ParseArguments(""),
				Router:    r,
			})
			return
		}

		// Check if message starts with one of defined prefixes.
		containsPrefix, content := hasPrefix(content, r.Prefixes, r.IgnorePrefixCase)
		if !containsPrefix {
			return
		}

		// Get rid of additional space.
		content = strings.TrimSpace(content)

		// if message is empty after prefix processing then do nothing.
		if content == "" {
			return
		}

		// split message at whitespace
		parts := SpliceRegex.Split(content, -1)

		for _, cmd := range r.Commands {
			if arrayContains(getIdentifiers(cmd), parts[0], cmd.IgnoreCase) {
				continue
			}
			content = strings.Join(parts[1:], " ")

			// define command context.
			ctx := &Context{
				Session:    session,
				Event:      event,
				Arguments:  ParseArguments(content),
				ObjectsMap: newObjectsMap(),
				Router:     r,
				Command:    cmd,
			}

			cmd.trigger(ctx)
		}
	}
}

func getIdentifiers(c *Command) []string {
	toCheck := make([]string, 0, len(c.Aliases)+1)
	toCheck = append(toCheck, c.Name)
	toCheck = append(toCheck, c.Aliases...)
	return toCheck
}
