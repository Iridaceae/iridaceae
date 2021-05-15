// Package components contains all static and support variables for iris.
package components

import (
	"errors"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	// ----------------------- build related -ldflags.

	AppVersion = "TEST_BUILD"
	AppCommit  = "TEST_BUILD"
	AppDate    = "0"
	Release    = "FALSE"
	Repo       = "https://github.com/Iridaceae/iridaceae"

	// ----------------------- statistic/metrics.

	StatsStartupTime      = time.Now()
	StatsCommandsExecuted = 0
	StatsMessageAnalyzed  = 0

	// ----------------------- errors definition.

	ErrSessionNotDefined = errors.New("session not defined")
	ErrEmbedNotDefined   = errors.New("embed not defined")

	// ----------------------- karma level.

	PermLvlBotOwner   = 69420
	PermLvlGuildOwner = 100

	DefaultAdminRules = []string{
		"+rs.guild.*",
		"+rs.etc.*",
		"+rs.chat.*",
	}
	DefaultUserRules = []string{
		"+rs.etc.*",
		"+rs.chat.*",
	}
	AdditionPermission = []string{
		"rs.guild.config.karma",
		"rs.guild.config.unbanrequest",
	}
)

const (

	// ----------------------- invitation permission.

	InvitePermission = 0x1 | // instant invite
		0x10 | // Manage channel
		0x20 | // manage guild
		0x40 | // add reaction
		0x400 | // view channel
		0x800 | // send messages
		0x2000 | // manage messages
		0x4000 | // embed links
		0x8000 | // attach files
		0x10000 | // read message history
		0x20000 | // mentions @everyone
		0x40000 | // use external emojis
		0x4000000 | // change nickname
		0x8000000 | // manage nickname
		0x10000000 | // manage roles
		0x20000000 | // manage webhooks
		0x40000000 // manage emoji

	// ----------------------- intent settings.

	Intents = discordgo.IntentsDirectMessages |
		discordgo.IntentsGuildBans |
		discordgo.IntentsGuildEmojis |
		discordgo.IntentsGuildIntegrations |
		discordgo.IntentsGuildInvites |
		discordgo.IntentsGuildMembers |
		discordgo.IntentsGuildMessageReactions |
		discordgo.IntentsGuildMessages |
		discordgo.IntentsGuildPresences |
		discordgo.IntentsGuildVoiceStates |
		discordgo.IntentsGuilds

	// ----------------------- colors.

	EmbedColorError    = 0xD32F2F
	EmbedColorDefault  = 0xFFC107
	EmbedColorUpdated  = 0x8BC34A
	EmbedColorGray     = 0xB0BEC5
	EmbedColorOrange   = 0xFB8C00
	EmbedColorGreen    = 0x8BC34A
	EmbedColorCyan     = 0x00BCD4
	EmbedColorYellow   = 0xFFC107
	EmbedColorViolet   = 0x6A1B9A
	ReportRevokedColor = 0x9C27B0
)

// IsRelease will check if given build tag is Release.
func IsRelease() bool {
	return Release == "TRUE"
}
