package jog

import (
	"errors"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/joho/godotenv"

	"github.com/stretchr/testify/assert"

	"github.com/Iridaceae/iridaceae/pkg"

	"github.com/bwmarrin/discordgo"
)

var (
	TestEmbed = &discordgo.MessageEmbed{
		Type:        "rich",
		Title:       "test pogu",
		Description: "good embed",
		Color:       0xffff00,
	}
)

func init() {
	// make sure to load env first
	err := godotenv.Load(strings.Join([]string{pkg.GetRootDir(), "defaults.env"}, "/"))
	if err != nil {
		panic(err)
	}
}

func makeTestSession() *discordgo.Session {
	var botToken = os.Getenv("CONCERTINA_AUTHTOKEN")
	// ensure sessions are established
	dg, err := discordgo.New("Bot " + botToken)
	if err != nil {
		panic(err)
	}
	return dg
}

func makeTestCtx() *Context {
	TestCtx := &Context{
		Session: &discordgo.Session{
			RWMutex: sync.RWMutex{},
			Token:   "test_token",
		},
		Event: &discordgo.MessageCreate{Message: &discordgo.Message{
			ID:        "test_msg",
			ChannelID: getEnvOrDefault("CONCERTINA_CHANNELID", ""),
			GuildID:   getEnvOrDefault("CONCERTINA_GUIDID", ""),
			Content:   "this is a test msg",
			Author: &discordgo.User{
				Username: "test_nick",
			},
			Embeds: []*discordgo.MessageEmbed{TestEmbed},
		}},
		Router:  TestRouter,
		Command: TestCommand,
	}
	return TestCtx
}

func TestContext_RespondText(t *testing.T) {
	ctx := makeTestCtx()
	ctx.Session = makeTestSession()
	err := ctx.RespondText("hello world")
	assert.Nil(t, err)
}

func TestContext_RespondEmbed(t *testing.T) {
	ctx := makeTestCtx()
	ctx.Session = makeTestSession()
	err := ctx.RespondEmbed(TestEmbed)
	assert.Nil(t, err)
}

func TestContext_RespondTextEmbed(t *testing.T) {
	ctx := makeTestCtx()
	ctx.Session = makeTestSession()
	err := ctx.RespondTextEmbed("Hello, this is a test with text and embed", TestEmbed)
	assert.Nil(t, err)
}

func TestContext_RespondTextEmbedError(t *testing.T) {
	ErrNothingWrong := errors.New("test error that will be printed to users about wrong cmd parser")
	ctx := makeTestCtx()
	ctx.Session = makeTestSession()
	err := ctx.RespondTextEmbedError("hello", "error response", ErrNothingWrong)
	assert.Nil(t, err)
}
