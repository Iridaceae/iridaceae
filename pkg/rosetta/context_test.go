package rosetta

import (
	"sync"
	"testing"
	"time"

	"github.com/Iridaceae/iridaceae/pkg"

	"github.com/stretchr/testify/assert"

	"github.com/Iridaceae/iridaceae/pkg/helpers"

	"github.com/sarulabs/di/v2"

	"github.com/bwmarrin/discordgo"
)

var (
	ctxAssert    = assert.New(&testing.T{})
	TestEmbedMsg = &discordgo.MessageEmbed{
		Type:        "rich",
		Title:       "This is a test message",
		Description: "Embed nice",
		Timestamp:   time.Now().Format(time.RFC3339),
		Color:       0xffff00,
	}
)

func init() {
	_ = pkg.LoadGlobalEnv()
}

func makeTestCtx(initOM bool) *context {
	ctx := &context{
		isDM:      true,
		isEdit:    true,
		args:      ParseArguments("a b c"),
		objectMap: &sync.Map{},
		session:   &discordgo.Session{Token: "rosetta_testToken"},
		message:   &discordgo.Message{ID: "rosetta_testMessage", Author: &discordgo.User{ID: "rosetta_testUser"}},
		guild:     &discordgo.Guild{ID: "rosetta_testGuild"},
		channel:   &discordgo.Channel{ID: helpers.GetEnvOrDefault("CONCERTINA_CHANNELID", "")},
		member:    &discordgo.Member{Nick: "rosetta_testNick"},
	}
	if initOM {
		b, _ := di.NewBuilder()
		err := b.Set("rosetta_testRouter", "rosetta_testValue")
		if err != nil {
			return nil
		}

		rr := &router{
			objectContainer: b.Build(),
			ctxPool:         &sync.Pool{New: func() interface{} { return &context{objectMap: &sync.Map{}} }},
			objectMap:       &sync.Map{},
		}
		ctx.router = rr
	}
	return ctx
}

func TestContextGetter(t *testing.T) {
	ctx := makeTestCtx(false)
	ctx.session = helpers.MakeTestSession()

	tests := []struct {
		name     string
		testFunc func() interface{}
		expected interface{}
	}{
		{"get session", func() interface{} { return ctx.GetSession() }, ctx.session},
		{"get arguments", func() interface{} { return ctx.GetArguments() }, ctx.args},
		{"get channel", func() interface{} { return ctx.GetChannel() }, ctx.channel},
		{"get message", func() interface{} { return ctx.GetMessage() }, ctx.message},
		{"get guild", func() interface{} { return ctx.GetGuild() }, ctx.guild},
		{"get user", func() interface{} { return ctx.GetUser() }, ctx.message.Author},
		{"get member", func() interface{} { return ctx.GetMember() }, ctx.member},
		{"is dm", func() interface{} { return ctx.IsDM() }, ctx.isDM},
		{"is edit", func() interface{} { return ctx.IsEdit() }, ctx.isEdit},
		{"respond text", func() interface{} {
			msg, err := ctx.RespondText("hello world")
			helpers.DeleteMessageAfter(ctx.GetSession(), msg, 20*time.Second)
			return err
		}, nil},
		{"respond embed", func() interface{} {
			msg, err := ctx.RespondEmbed(TestEmbedMsg)
			defer helpers.DeleteMessageAfter(ctx.GetSession(), msg, 20*time.Second)
			return err
		}, nil},
		{"respond embed error", func() interface{} {
			msg, err := ctx.RespondEmbedError("test with defined error", ErrCommandNotFound)
			defer helpers.DeleteMessageAfter(ctx.GetSession(), msg, 20*time.Second)
			return err
		}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctxAssert.Equal(tt.expected, tt.testFunc())
		})
	}
}

func TestContext_SetGetInitObjectMap(t *testing.T) {
	t.Run("get set test local", func(t *testing.T) {
		cc := makeTestCtx(true)
		cc.SetObject("key", 123)
		v, ok := cc.GetObject("key").(int)
		ctxAssert.True(ok)
		ctxAssert.Equal(123, v)

		v, ok = cc.GetObject("unexisted key").(int)
		ctxAssert.False(ok)
		ctxAssert.Equal(0, v)
	})
	t.Run("get om global", func(t *testing.T) {
		ctx := makeTestCtx(true)
		v, ok := ctx.GetObject("rosetta_testRouter").(string)
		if !ok {
			t.Error("recovered global value should have type string")
		}
		if v != "rosetta_testValue" {
			t.Error("invalid global value")
		}
	})
	t.Run("init om", func(t *testing.T) {

	})
}
