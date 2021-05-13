// Package helpers contains helpers functions that can be use globally.
package helpers

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

// DeleteMessageAfter tries to delete given message after a specified duration.
func DeleteMessageAfter(s *discordgo.Session, msg *discordgo.Message, duration time.Duration) {
	if msg == nil {
		return
	}
	time.AfterFunc(duration, func() {
		_ = s.ChannelMessageDelete(msg.ChannelID, msg.ID)
	})
}
