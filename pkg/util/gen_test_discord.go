package util

import (
	"os"

	"github.com/bwmarrin/discordgo"
)

// TODO: a better way to get general token.
func MakeTestSession() *discordgo.Session {
	botToken := os.Getenv("CONCERTINA_AUTHTOKEN")
	// ensure sessions are established
	dg, err := discordgo.New("Bot " + botToken)
	if err != nil {
		panic(err)
	}
	return dg
}
