package rosetta

import (
	"time"
)

const (
	// GroupGlobalAdmin defines stringer for global admin.
	GroupGlobalAdmin = "GLOBAL ADMIN"
	GroupGuildAdmin  = "GUILD ADMIN"
	GroupModeration  = "MODERATION"
	GroupFun         = "FUN"
	GroupChat        = "CHAT"
	GroupEtc         = "ETC"
	GroupGeneral     = "GENERAL"
	GroupGuildConfig = "GUILD CONFIG"
)

// LimitedConfig defines command that is rate limit-able.
type LimitedConfig interface {

	// GetLimiterBurst returns max amount of tokens which can be available at time.
	GetLimiterBurst() int

	// GetLimiterRestoration returns duration between new token get generated.
	GetLimiterRestoration() time.Duration

	// IsLimiterGlobal return true if limit shall be handled globally across all guilds.
	// Otherwise it should be created independently.
	IsLimiterGlobal() bool
}

// SubPermission wraps information about a command sub permission.
type SubPermission struct {
	Term        string `json:"term"`
	Explicit    bool   `json:"explicit"`
	Description string `json:"description"`
}

// Command defines a functionality of a command struct that will be registered under router.
type Command interface {
	// GetInvokers defines a unique string defines command invokers. First will be the primary
	// commands, following with aliases.
	GetInvokers() []string

	// GetDescription describes how the function will work.
	GetDescription() string

	// GetUsage returns how one can use the command with its subcommands.
	GetUsage() string

	// GetGroup returns the groups command belongs to.
	// admin - user - bot - helpers - etc.
	GetGroup() string

	// GetDomain returns the commands domain name.
	// TODO: have a regex to check if it follows the correct domain definition.
	// Domain definition as follow: rs.group.main.etc ...
	GetDomain() string

	// GetSubPermissionRules returns optional sub permissions of command.
	GetSubPermissionRules() []SubPermission

	// IsExecutableInDM returns true when command can be used when user dms the bot.
	// I saw this in yagpdb and I would like for rosetta to have something similar.
	IsExecutableInDM() bool

	// Exec is called when command is executed and getting passed CommandArgs.
	// Returns nil when successfully executed, otherwise errors encountered will be returned.
	Exec(ctx Context) error
}
