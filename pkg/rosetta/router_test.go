package rosetta

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Iridaceae/iridaceae/internal/helpers"

	"github.com/bwmarrin/discordgo"

	"github.com/rs/zerolog/log"

	"github.com/stretchr/testify/assert"
)

var TestSession *discordgo.Session

func init() {
	_ = helpers.LoadGlobalEnv()
	TestSession = helpers.MakeTestSession()
}

func TestRouter_Setup(t *testing.T) {
	r := NewRouter(makeTestConfig())
	r.RegisterSession(TestSession)
}

func TestNewDefaultConfig(t *testing.T) {
	ctx := makeTestCtx(false, true)
	c := NewRouterConfig()
	c.OnError(ctx, ErrTypeMiddleware, ErrMiddleware)
}

func TestNewRouter(t *testing.T) {
	ctx := makeTestCtx(false, true)
	tests := []struct {
		cfg                *RouterConfig
		cOnError           bool
		cGuildPrefixGetter bool
	}{
		{NewRouterConfig(), false, false},
		{makeTestConfig(), true, true},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d-routertest", i), func(t *testing.T) {
			if !tt.cOnError {
				tt.cfg.OnError = nil
			}
			if !tt.cGuildPrefixGetter {
				tt.cfg.GuildPrefixGetter = nil
			}
			r := NewRouter(tt.cfg)
			if !tt.cOnError {
				r.GetConfig().OnError(ctx, ErrTypeGetGuild, ErrGetGuild)
			}
			if !tt.cGuildPrefixGetter {
				out, err := r.GetConfig().GuildPrefixGetter(ctx.GetGuild().ID)
				assert.Empty(t, out)
				assert.Nil(t, err)
			}
			if tt.cfg.UseDefaultHelpCommand {
				// we have 4 aliases for help command hence length = 4
				assert.Len(t, r.GetCommandMap(), 4)
				assert.Len(t, r.GetCommandInstances(), 1)
			}
		})
	}
}

func TestRouter_Register(t *testing.T) {
	t.Run("register command", func(t *testing.T) {
		r := NewRouter(makeTestConfig())
		cmd := &TestCmd{}
		r.Register(cmd)

		for _, instance := range r.(*routerImpl).cmdMap {
			assert.Equal(t, cmd, instance)
		}
	})
	t.Run("panic register command", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("register already invoked command didn't panic")
			}
		}()
		r := NewRouter(makeTestConfig())
		cmd := &TestCmd{}
		r.Register(cmd)
		r.Register(cmd)
	})
	t.Run("register middleware", func(t *testing.T) {
		r := NewRouter(makeTestConfig())
		mw := &TestMiddleware{}
		r.Register(mw)
		assert.NotZero(t, r.(*routerImpl).middleware)
		assert.Equal(t, mw, r.(*routerImpl).middleware[0])
	})
	t.Run("register panic interface", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("register invalid instance didn't panic")
			}
		}()
		r := NewRouter(makeTestConfig())
		r.Register("invalid")
	})
}

func TestRouterExecuteMiddleware(t *testing.T) {
	// this will determine how many times we want to run our test channel.
	const count = 2
	exit := make(chan bool, count)
	session := helpers.MakeTestSession()

	cmd := &TestCmd{}
	cfg := makeTestConfig()
	cfg.AllowBots = true

	r := NewRouter(cfg)
	r.Register(cmd)

	m1 := &TestMiddleware{}
	m1.layer = LayerBeforeCommand
	r.Register(m1)
	m2 := &TestMiddleware{}
	m2.layer = LayerAfterCommand
	r.Register(m2)

	msg := makeTestMsg(t, "!ping")

	session.AddHandler(func(s *discordgo.Session, e *discordgo.Ready) {
		failFunc := func(s *discordgo.Session, fail bool, msg *discordgo.Message) {
			cmd.fail = fail
			r.(*routerImpl).trigger(s, msg)
			switch cmd.fail {
			case true:
				assert.True(t, m1.executed)
				assert.False(t, m2.executed)
			case false:
				assert.True(t, m1.executed)
				assert.True(t, m1.executed)
			default:
				assert.False(t, m1.executed)
				assert.False(t, m2.executed)
			}
			exit <- true
		}
		failFunc(s, true, msg)
		failFunc(s, false, msg)
	})
	if err := session.Open(); err != nil {
		log.Err(err).Msg("")
	}
	<-exit
}

func TestRouterExecuteCommand(t *testing.T) {
	exit := make(chan bool, 4)

	tests := []struct {
		name          string
		shallExecuted bool
		testHandler   func(m *discordgo.Message)
	}{
		{"author is bot", false, func(m *discordgo.Message) {
			m.Author.Bot = true
			m.Content = "!ping"
		}},
		{"use default prefix to execute", true, func(m *discordgo.Message) { m.Content = "!ping" }},
		{"will execute with custom prefix", true, func(m *discordgo.Message) { m.Content = "test!ping" }},
		{"command does not exist", false, func(m *discordgo.Message) { m.Content = "!abc" }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testTrigger(t, TestSession, exit, tt.shallExecuted, tt.testHandler)
		})
	}

	if err := TestSession.Open(); err != nil {
		log.Err(err).Msg("")
	}
	<-exit
}

// NOTE: current tests is running very slow. Potential qol may be to have a general session?
func testTrigger(t *testing.T, session *discordgo.Session, exit chan bool, shallExecuted bool, configurator func(m *discordgo.Message)) {
	t.Helper()

	cmd := &TestCmd{}
	cfg := makeTestConfig()
	cfg.AllowBots = true
	cfg.DeleteMessageAfter = false

	r := NewRouter(cfg)
	r.Register(cmd)

	msg := makeTestMsg(t)
	configurator(msg)
	r.GetConfig().OnError = func(ctx Context, errType ErrorType, err error) {}

	session.AddHandler(func(_ *discordgo.Session, e *discordgo.Ready) {
		r.(*routerImpl).trigger(session, msg)
		if !cmd.executed && shallExecuted {
			t.Error("command was not executed")
		} else if cmd.executed && !shallExecuted {
			t.Error("command was executed")
		}
		exit <- true
	})
}

func TestRouter_GetterSetter(t *testing.T) {
	cfg := makeTestConfig()
	cmd := &TestCmd{}
	r := NewRouter(cfg)
	r.Register(cmd)

	tests := []struct {
		name     string
		testFunc func() interface{}
		expected interface{}
	}{
		{"get configparser", func() interface{} { return r.GetConfig() }, cfg},
		{"get command map", func() interface{} { return r.GetCommandMap() }, map[string]Command{"ping": &TestCmd{}, "p": &TestCmd{}}},
		{"get command instance", func() interface{} { return r.GetCommandInstances()[0] }, cmd},
		{"get command", func() interface{} {
			res := make([]bool, 0)
			for i := range r.(*routerImpl).cmdMap {
				_, ok := r.GetCommand(i)
				res = append(res, ok)
			}
			return res
		}, []bool{true, true}},
		{"set object", func() interface{} {
			r.SetObject("rosetta_testAnother", 69420)
			rec, ok := r.(*routerImpl).objectMap.Load("rosetta_testAnother")
			assert.True(t, ok)
			v, _ := rec.(int)
			return v
		}, 69420},
		{"get object", func() interface{} {
			r.(*routerImpl).objectMap.Store("rosetta_test", 456)
			v, _ := r.GetObject("rosetta_test").(int)
			return v
		}, 456},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.testFunc())
		})
	}
}

func makeTestMsg(t *testing.T, content ...string) *discordgo.Message {
	t.Helper()
	cmd := &discordgo.Message{
		ChannelID: helpers.GetEnvOrDefault("CONCERTINA_CHANNELID", ""),
		GuildID:   helpers.GetEnvOrDefault("CONCERTINA_GUILDID", ""),
		Author:    &discordgo.User{ID: helpers.GetEnvOrDefault("CONCERTINA_USERID", ""), Bot: false},
		Member: &discordgo.Member{
			GuildID: helpers.GetEnvOrDefault("CONCERTINA_GUILDID", ""),
			User:    &discordgo.User{ID: helpers.GetEnvOrDefault("CONCERTINA_USERID", ""), Bot: false},
		},
	}
	if len(content) > 0 && content[0] != "" {
		cmd.Content = content[0]
	}
	return cmd
}

func makeTestConfig() *RouterConfig {
	return &RouterConfig{
		GeneralPrefix:         "!",
		IgnoreCase:            true,
		AllowDM:               false,
		AllowBots:             false,
		ExecuteOnEdit:         false,
		UseDefaultHelpCommand: false,
		DeleteMessageAfter:    false,
		OnError:               func(_ Context, errType ErrorType, err error) { log.Err(err).Msgf("type [%d]", errType) },
		GuildPrefixGetter:     func(string) (string, error) { return "test!", nil },
	}
}

type TestCmd struct {
	executed bool
	fail     bool
}

func (t *TestCmd) GetInvokers() []string {
	return []string{"ping", "p"}
}

func (t *TestCmd) GetDescription() string {
	return "ping pong ding dong"
}

func (t *TestCmd) GetUsage() string {
	return "`ping` - ping"
}

func (t *TestCmd) GetGroup() string {
	return GroupFun
}

func (t *TestCmd) GetDomain() string {
	return "test.fun.ping"
}

func (t *TestCmd) GetSubPermissionRules() []SubPermission {
	return nil
}

func (t *TestCmd) IsExecutableInDM() bool {
	return true
}

func (t *TestCmd) Exec(_ Context) error {
	t.executed = true
	if t.fail {
		return errors.New("test error")
	}
	return nil
}
