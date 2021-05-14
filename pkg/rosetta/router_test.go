package rosetta

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Iridaceae/iridaceae/pkg"

	"github.com/bwmarrin/discordgo"

	"github.com/rs/zerolog/log"

	"github.com/stretchr/testify/assert"

	"github.com/Iridaceae/iridaceae/pkg/helpers"
)

func init() {
	_ = pkg.LoadGlobalEnv()
}

func TestRouter_Setup(t *testing.T) {
	r := NewRouter(makeTestConfig())
	s := helpers.MakeTestSession()
	r.Setup(s)
}

func TestNewDefaultConfig(t *testing.T) {
	ctx := makeTestCtx(false)
	ctx.session = helpers.MakeTestSession()

	c := NewDefaultConfig()
	c.OnError(ctx, ErrTypeMiddleware, ErrMiddleware)
}

func TestNewRouter(t *testing.T) {
	ctx := makeTestCtx(true)
	ctx.session = helpers.MakeTestSession()
	tests := []struct {
		cfg                *Config
		cOnError           bool
		cGuildPrefixGetter bool
	}{
		{NewDefaultConfig(), false, false},
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

		for _, instance := range r.(*router).cmdMap {
			assert.Equal(t, cmd, instance)
		}
	})
	t.Run("register middleware", func(t *testing.T) {
		r := NewRouter(makeTestConfig())
		mw := &TestMiddleware{}
		r.Register(mw)
		assert.NotZero(t, r.(*router).middleware)
		assert.Equal(t, mw, r.(*router).middleware[0])
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
	const count = 3
	session := helpers.MakeTestSession()
	exit := make(chan bool, count)

	cmd := &TestCmd{}
	cfg := makeTestConfig()
	cfg.AllowBots = true
	cfg.DeleteMessageAfter = false

	r := NewRouter(cfg)
	r.Register(cmd)

	m1 := &TestMiddleware{}
	m1.layer = LayerBeforeCommand
	r.Register(m1)
	m2 := &TestMiddleware{}
	m2.layer = LayerAfterCommand
	r.Register(m2)

	var msg = makeTestMsg(t, "!ping")

	session.AddHandler(func(_ *discordgo.Session, e *discordgo.Ready) {
		failFunc := func(fail bool, msg *discordgo.Message) {
			cmd.fail = fail
			r.(*router).trigger(session, msg)
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

		failFunc(true, msg)
		failFunc(false, msg)
		msg.Content = ""
		failFunc(false, msg)
	})
	_ = session.Open()
	<-exit
}

func TestRouterExecuteCommand(t *testing.T) {
	tests := []struct {
		name          string
		shallExecuted bool
		testHandler   func(m *discordgo.Message)
	}{
		{"fail to execute", true, func(m *discordgo.Message) { m.Content = "!ping" }},
		{"will execute", true, func(m *discordgo.Message) { m.Content = "test!ping" }},
		{"command does not exist", false, func(m *discordgo.Message) { m.Content = "!abc" }},
		{"author is bot", false, func(m *discordgo.Message) {
			m.Author.Bot = true
			m.Content = "!ping"
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			testTrigger(t, tt.shallExecuted, tt.testHandler)
		})
	}
}

// NOTE: current tests is running very slow. Potential qol may be to have a general session?
func testTrigger(t *testing.T, shallExecuted bool, configurator func(m *discordgo.Message)) {
	t.Helper()

	session := helpers.MakeTestSession()
	exit := make(chan bool, 1)

	cmd := &TestCmd{}
	cfg := makeTestConfig()
	cfg.AllowBots = true
	cfg.DeleteMessageAfter = false

	r := NewRouter(cfg)
	r.Register(cmd)

	msg := makeTestMsg(t)
	configurator(msg)

	session.AddHandler(func(_ *discordgo.Session, e *discordgo.Ready) {
		r.(*router).trigger(session, msg)
		if !cmd.executed && shallExecuted {
			t.Error("command was not executed")
		} else if cmd.executed && !shallExecuted {
			t.Error("command was executed")
		}
		exit <- true
	})

	err := session.Open()
	if err != nil {
		log.Err(err).Msg("")
	}
	<-exit
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
		{"get config", func() interface{} { return r.GetConfig() }, cfg},
		{"get command map", func() interface{} { return r.GetCommandMap() }, map[string]Command{"ping": &TestCmd{}, "p": &TestCmd{}}},
		{"get command instance", func() interface{} { return r.GetCommandInstances()[0] }, cmd},
		{"get command", func() interface{} {
			res := make([]bool, 0)
			for i := range r.(*router).cmdMap {
				_, ok := r.GetCommand(i)
				res = append(res, ok)
			}
			return res
		}, []bool{true, true}},
		{"set object", func() interface{} {
			r.SetObject("rosetta_testAnother", 69420)
			rec, ok := r.(*router).objectMap.Load("rosetta_testAnother")
			assert.True(t, ok)
			v, _ := rec.(int)
			return v
		}, 69420},
		{"get object", func() interface{} {
			r.(*router).objectMap.Store("rosetta_test", 456)
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

func (t *TestCmd) Exec(ctx Context) error {
	t.executed = true
	if t.fail {
		return errors.New("test error")
	}
	log.Debug().Msg(ctx.GetMessage().Content)
	return nil
}

func makeTestConfig() *Config {
	return &Config{
		GeneralPrefix:         "!",
		IgnoreCase:            true,
		AllowDM:               false,
		AllowBots:             false,
		ExecuteOnEdit:         false,
		UseDefaultHelpCommand: false,
		DeleteMessageAfter:    true,
		OnError:               func(_ Context, errType ErrorType, err error) { log.Err(err).Msgf("type [%d]", errType) },
		GuildPrefixGetter:     func(string) (string, error) { return "test!", nil },
	}
}
