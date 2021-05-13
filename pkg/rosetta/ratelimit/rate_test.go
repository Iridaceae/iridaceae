package ratelimit

import (
	"testing"

	"github.com/bwmarrin/discordgo"

	"github.com/stretchr/testify/assert"

	"github.com/Iridaceae/iridaceae/pkg/rosetta"
)

func testLoop(t *testing.T, m ...Manager) *RateLimiter {
	t.Helper()

	tm := New(m...)
	cmd := &TestCmd{false, false, false}
	ctx := &TestContext{
		chanType: discordgo.ChannelTypeGuildText,
		gid:      "gid",
		uid:      "uid",
	}

	pass := func() {
		ok, err := tm.Handle(cmd, ctx, tm.GetLayer())
		assert.Nil(t, err)
		assert.True(t, ok, "rate limiter stopped unexpectedly")
	}
	fail := func() {
		ok, err := tm.Handle(cmd, ctx, tm.GetLayer())
		assert.Nil(t, err)
		assert.False(t, ok, "rate limiter passed unexpectedly")
	}

	for i := 0; i < cmd.GetLimiterBurst(); i++ {
		pass()
	}
	fail()

	return tm
}

func TestRateLimiter_GetLayer(t *testing.T) {
	tm := testLoop(t)
	assert.Equal(t, rosetta.LayerBeforeCommand, tm.GetLayer())
}

// TODO: custom manager test.
func TestRateLimiter_Handle(t *testing.T) {
	t.Run("test rate limiter with default manager", func(t *testing.T) {
		testLoop(t)
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
