package rosetta

import (
	"sync"
	"testing"

	"github.com/Iridaceae/iridaceae/pkg/util"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

func TestHasPrefix(t *testing.T) {
	t.Run("doesn't have prefix contain in given string", func(t *testing.T) {
		s := "hello world"
		prefs := []string{"!", "-"}
		ok, _ := hasPrefix(s, prefs, true)
		assert.False(t, ok)
	})
	t.Run("does have prefix in given string", func(t *testing.T) {
		s := "!hello world"
		prefs := []string{"!", "-"}
		ok, so := hasPrefix(s, prefs, true)
		assert.True(t, ok)
		assert.Equal(t, "hello world", so)
	})
}

func TestTrimPreSuffix(t *testing.T) {
	s := "'hello world'"
	preSuffix := "'"
	o := trimPreSuffix(s, preSuffix)
	assert.Equal(t, "hello world", o)
}

func TestArrayContains(t *testing.T) {
	tarr := []string{"1", "2", "3"}
	contained := "test"
	ok := arrayContains(tarr, contained, false)
	assert.False(t, ok)
}

func makeTestCtx() *Context {
	testCtx := &Context{
		Session: &discordgo.Session{
			RWMutex: sync.RWMutex{},
			Token:   "test_token",
		},
		Arguments: &Arguments{
			raw: "test t1 t2",
		},
		Event: &discordgo.MessageCreate{Message: &discordgo.Message{
			ID:        "test_msg",
			ChannelID: util.GetEnvOrDefault("CONCERTINA_CHANNELID", ""),
			GuildID:   util.GetEnvOrDefault("CONCERTINA_GUILDID", ""),
			Content:   "this is a test msg",
			Author: &discordgo.User{
				ID:       "12341234",
				Username: "test_nick",
			},
			Embeds: []*discordgo.MessageEmbed{TestEmbedMsg},
		}},
		Router:  TestRouter,
		Command: TestCommand,
	}
	return testCtx
}
