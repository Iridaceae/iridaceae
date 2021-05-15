// Package helpers contains helpers functions that can be use globally.
package helpers

import (
	"fmt"
	"os"
	"time"

	"github.com/Iridaceae/iridaceae/pkg/configmanager"

	"github.com/bwmarrin/discordgo"
)

// MakeTestSession returns a discordgo.Session for testing.
func MakeTestSession() *discordgo.Session {
	// TODO: a better way to get general token.
	botToken := os.Getenv("CONCERTINA_AUTHTOKEN")
	// ensure sessions are established
	dg, err := discordgo.New("Bot " + botToken)
	if err != nil {
		panic(err)
	}
	return dg
}

// DeleteMessageAfter tries to delete given message after a specified duration.
func DeleteMessageAfter(s *discordgo.Session, msg *discordgo.Message, duration time.Duration) {
	if msg == nil {
		return
	}
	time.AfterFunc(duration, func() {
		_ = s.ChannelMessageDelete(msg.ChannelID, msg.ID)
	})
}

// GetInviteLink returns bot invite link with correct permissions.
func GetInviteLink(cid configmanager.Options) string {
	return fmt.Sprintf(BaseAuthURLTemplate, cid.ToString(), InvitePermission)
}
