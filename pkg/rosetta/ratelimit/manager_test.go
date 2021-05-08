package ratelimit

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/joho/godotenv"

	"github.com/Iridaceae/iridaceae/pkg"

	"github.com/stretchr/testify/assert"

	"github.com/bwmarrin/discordgo"

	"github.com/Iridaceae/iridaceae/pkg/rosetta"
)

func onTestCmd(ctx *rosetta.Context, _ ...interface{}) {
	_ = ctx.RespondText(strconv.Itoa(ctx.ObjectsMap.GetValue("testObject").(int)))
}

var (
	Embed = &discordgo.MessageEmbed{
		Type:        "rich",
		Title:       "This is a test message",
		Description: "Embed nice",
		Timestamp:   time.Now().Format(time.RFC3339),
		Color:       0xffff00,
	}

	Router = &rosetta.Router{
		Prefixes:         []string{"!"},
		IgnorePrefixCase: false,
		BotsAllowed:      false,
		Commands:         []*rosetta.Command{},
		Middlewares:      []rosetta.Middleware{},
		PingHandler:      func(context *rosetta.Context, _ ...interface{}) { _ = context.RespondText("pong!") },
	}

	Command = &rosetta.Command{
		Name:        "test",
		Aliases:     []string{"test", "t"},
		Description: "test commands",
		IgnoreCase:  true,
		SubCommands: []*rosetta.Command{},
		Handler:     onTestCmd,
	}

	TestCtx = &rosetta.Context{
		Session: &discordgo.Session{
			RWMutex: sync.RWMutex{},
			Token:   "test_token",
		},
		Event: &discordgo.MessageCreate{Message: &discordgo.Message{
			ID:        "test_msg",
			ChannelID: os.Getenv("CONCERTINA_CHANNELID"),
			GuildID:   os.Getenv("CONCERTINA_GUILDID"),
			Content:   "this is a test msg",
			Author: &discordgo.User{
				ID:       "12341234",
				Username: "test_nick",
			},
			Embeds: []*discordgo.MessageEmbed{Embed},
		}},
		Router:  Router,
		Command: Command,
	}

	TestDiffCtx = &rosetta.Context{
		Session: &discordgo.Session{
			RWMutex: sync.RWMutex{},
			Token:   "test_token",
		},
		Event: &discordgo.MessageCreate{Message: &discordgo.Message{
			ID:        "different_test_msg",
			ChannelID: "768165783062052884",
			GuildID:   "723280184257544193",
			Content:   "another test message",
			Author: &discordgo.User{
				ID:       "123123123",
				Username: "test_nick_2",
			},
			Embeds: []*discordgo.MessageEmbed{Embed},
		}},
		Router:  Router,
		Command: Command,
	}
)

func init() {
	// make sure to load env first
	err := godotenv.Load(strings.Join([]string{pkg.GetRootDir(), "defaults.env"}, "/"))
	if err != nil {
		return
	}
}

func TestManager_GetBucket(t *testing.T) {
	m := newManager(10 * time.Minute)

	l1 := m.GetBucket(TestCtx)
	l2 := m.GetBucket(TestCtx)
	l3 := m.GetBucket(TestDiffCtx)
	assert.Equal(t, l1, l2)
	if l3 == l1 || l3 == l2 {
		t.Error(errDupsBucket(l3))
	}
}

func errDupsBucket(i1 interface{}) string {
	return fmt.Sprintf("%+v was a duplicate of l1 & l2.", i1)
}
