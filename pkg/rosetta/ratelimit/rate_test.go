package ratelimit

import (
	"fmt"
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/stretchr/testify/assert"

	"github.com/Iridaceae/iridaceae/pkg/rosetta"
)

func testLoop(t *testing.T, m ...Manager) *RateLimiter {
	t.Helper()

	rl := New(m...)
	cmd := &TestCmd{false, false, false}
	ctx := &TestContext{
		chanType: discordgo.ChannelTypeGuildText,
		gid:      "gid",
		uid:      "uid",
	}

	pass := func() {
		ok, err := rl.Handle(cmd, ctx, rl.GetLayer())
		assert.Nil(t, err)
		assert.True(t, ok, "rate limiter stopped unexpectedly")
	}
	fail := func() {
		ok, err := rl.Handle(cmd, ctx, rl.GetLayer())
		assert.Nil(t, err)
		assert.False(t, ok, "rate limiter passed unexpectedly")
	}

	for i := 0; i < cmd.GetLimiterBurst(); i++ {
		pass()
	}
	fail()

	return rl
}

func TestRateLimiter_GetLayer(t *testing.T) {
	rl := testLoop(t)
	assert.Equal(t, rosetta.LayerBeforeCommand, rl.GetLayer())
}

// TODO: custom manager test.
func TestRateLimiter_Handle(t *testing.T) {
	t.Run("test rate limiter with default manager", func(t *testing.T) {
		testLoop(t)
	})
	t.Run("test rate limiter with multiple custom managers", func(t *testing.T) {
		cm := newInternalManager(20 * time.Minute)
		rl := testLoop(t, cm, newInternalManager(30*time.Minute))
		assert.Equal(t, rl.m, cm)
	})
	t.Run("handled a non-implemented commands", func(t *testing.T) {
		rl := New()
		cmd := &TestCmdNotImplemented{}
		ctx := &TestContext{
			chanType: discordgo.ChannelTypeGuildText,
			gid:      "gid",
			uid:      "uid",
		}
		ok, err := rl.Handle(cmd, ctx, rl.GetLayer())
		assert.True(t, ok)
		assert.Nil(t, err)
	})
	t.Run("global dms", func(t *testing.T) {
		rl := New()
		cmd := &TestCmd{false, false, true}
		ctx := &TestContext{
			chanType: discordgo.ChannelTypeDM,
			gid:      "gid",
			uid:      "uid",
		}

		expected := fmt.Sprintf("%s:%s:%s", cmd.GetDomain(), ctx.GetUser().ID, "__global__")
		_, _ = rl.Handle(cmd, ctx, rl.GetLayer())
		_, ok := rl.m.GetExecutions().GetValue(expected).(*Bucket)
		assert.True(t, ok)
	})
	t.Run("in the dms", func(t *testing.T) {
		rl := New()
		cmd := &TestCmd{false, false, false}
		ctx := &TestContext{
			chanType: discordgo.ChannelTypeDM,
			gid:      "gid",
			uid:      "uid",
		}

		expected := fmt.Sprintf("%s:%s:%s", cmd.GetDomain(), ctx.GetUser().ID, "__dm__")
		_, _ = rl.Handle(cmd, ctx, rl.GetLayer())
		_, ok := rl.m.GetExecutions().GetValue(expected).(*Bucket)
		assert.True(t, ok)
	})
}

type TestContext struct {
	chanType discordgo.ChannelType
	gid      string
	uid      string
}

func (tc *TestContext) GetObject(key string) (value interface{}) {
	return nil
}

func (tc *TestContext) SetObject(key string, value interface{}) {}

func (tc *TestContext) GetSession() *discordgo.Session {
	return nil
}

func (tc *TestContext) GetArguments() *rosetta.Arguments {
	return nil
}

func (tc *TestContext) GetChannel() *discordgo.Channel {
	return &discordgo.Channel{Type: tc.chanType}
}

func (tc *TestContext) GetMessage() *discordgo.Message {
	return nil
}

func (tc *TestContext) GetGuild() *discordgo.Guild {
	return &discordgo.Guild{ID: tc.gid}
}

func (tc *TestContext) GetUser() *discordgo.User {
	return &discordgo.User{ID: tc.uid}
}

func (tc *TestContext) GetMember() *discordgo.Member {
	return nil
}

func (tc *TestContext) IsDM() bool {
	return false
}

func (tc *TestContext) IsEdit() bool {
	return false
}

func (tc *TestContext) RespondText(content string) (*discordgo.Message, error) {
	return nil, nil
}

func (tc *TestContext) RespondEmbed(embed *discordgo.MessageEmbed) (*discordgo.Message, error) {
	return nil, nil
}

func (tc *TestContext) RespondEmbedError(title string, err error) (*discordgo.Message, error) {
	return nil, nil
}

func (tc *TestContext) RespondTextEmbed(content string, embed *discordgo.MessageEmbed) (*discordgo.Message, error) {
	return nil, nil
}

func (tc *TestContext) RespondTextEmbedError(title, content string, err error) (*discordgo.Message, error) {
	return nil, nil
}

type TestCmdNotImplemented struct{}

func (t *TestCmdNotImplemented) GetInvokers() []string {
	return nil
}

func (t *TestCmdNotImplemented) GetDescription() string {
	return ""
}

func (t *TestCmdNotImplemented) GetUsage() string {
	return ""
}

func (t *TestCmdNotImplemented) GetGroup() string {
	return ""
}

func (t *TestCmdNotImplemented) GetDomain() string {
	return ""
}

func (t *TestCmdNotImplemented) GetSubPermissionRules() []rosetta.SubPermission {
	return nil
}

func (t *TestCmdNotImplemented) IsExecutableInDM() bool {
	return false
}

func (t *TestCmdNotImplemented) Exec(ctx rosetta.Context) error {
	return rosetta.ErrCommandExec
}
