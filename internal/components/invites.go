package components

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

const (
	BaseAuthURLTemplate string = "https://discord.com/api/oauth2/authorize?client_id=%s&scope=bot&permissions=%d"
)

// GetInviteLink returns bot invite link with correct permissions.
func GetInviteLink(s *discordgo.Session) string {
	return fmt.Sprintf(BaseAuthURLTemplate, s.State.User.ID, InvitePermission)
}
