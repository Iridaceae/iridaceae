package commands

import (
	"github.com/bwmarrin/discordgo"
)

// Route holds information about a specific message route handler.
// This will be responsible to routing input user message to specific commands.
type Route struct {
	Pattern     string
	Description string
	Help        string
	Runner      HandlerFunc
}

// Context holds a bit of extra data rather than current string implementation.
// As we move on for more generic command parsing for Iris, it should take input from users as a Context rather than single string parsing.
// By using this we only have to process some of the info once.
type Context struct {
	Fields          []string
	Content         string
	IsDirected      bool
	IsPrivate       bool
	HasPrefix       bool
	HasMention      bool
	HasMentionFirst bool
}

type HandlerFunc func(s *discordgo.Session, m *discordgo.Message, ctx *Context)
