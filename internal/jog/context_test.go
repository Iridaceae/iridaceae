package jog

import (
	"errors"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/joho/godotenv"

	"github.com/stretchr/testify/assert"

	"github.com/Iridaceae/iridaceae/pkg"

	"github.com/bwmarrin/discordgo"
)

var TestEmbed = &discordgo.MessageEmbed{
	Type:        "rich",
	Title:       "This is a test message",
	Description: "Embed nice",
	Timestamp:   time.Now().Format(time.RFC3339),
	Color:       0xffff00,
}

var TestDiffCtx = &Context{
	Session: &discordgo.Session{
		RWMutex: sync.RWMutex{},
		Token:   "test_token",
	},
	Event: &discordgo.MessageCreate{Message: &discordgo.Message{
		ID:        "different_test_msg",
		ChannelID: "840228494009171988",
		GuildID:   "832847806162403358",
		Content:   "another test message",
		Author: &discordgo.User{
			ID:       "123123123",
			Username: "test_nick_2",
		},
		Embeds: []*discordgo.MessageEmbed{TestEmbed},
	}},
	Router:  TestRouter,
	Command: TestCommand,
}

func init() {
	// make sure to load env first
	err := godotenv.Load(strings.Join([]string{pkg.GetRootDir(), "defaults.env"}, "/"))
	if err != nil {
		TestLogger.Warn("defaults.env not found. This is due to either Docker container or CI is running the tasks. Loading from ENVARS instead.")
	}
}

func makeTestSession() *discordgo.Session {
	botToken := os.Getenv("CONCERTINA_AUTHTOKEN")
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
			GuildID:   getEnvOrDefault("CONCERTINA_GUILDID", ""),
			Content:   "this is a test msg",
			Author: &discordgo.User{
				ID:       "12341234",
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
	err := ctx.RespondTextEmbedError("Hello, this is a error text embed test", "error response", ErrNothingWrong)
	assert.Nil(t, err)
}
