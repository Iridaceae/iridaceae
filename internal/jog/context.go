package jog

import (
	"github.com/bwmarrin/discordgo"
)

// Context represents the context for a command event.
type Context struct {
	Session    *discordgo.Session
	Event      *discordgo.MessageCreate
	Arguments  *Arguments
	ObjectsMap *ObjectsMap
	Router     *Router
	Command    *Command
}

// ExecutionHandler represents a handler for a context execution.
type ExecutionHandler func(*Context)

// RespondText wraps around responses of given text message.
func (c *Context) RespondText(text string) error {
	_, err := c.Session.ChannelMessageSend(c.Event.ChannelID, text)
	return err
}

// RespondEmbed responds with the given embed message.
func (c *Context) RespondEmbed(embed *discordgo.MessageEmbed) error {
	_, err := c.Session.ChannelMessageSendEmbed(c.Event.ChannelID, embed)
	return err
}

// RespondTextEmbed responds with given text and embed message.
func (c *Context) RespondTextEmbed(text string, embed *discordgo.MessageEmbed) error {
	_, err := c.Session.ChannelMessageSendComplex(c.Event.ChannelID, &discordgo.MessageSend{
		Content: text,
		Embed:   embed,
	})
	return err
}

// RespondTextEmbedError responds given error to users with embed message.
func (c *Context) RespondTextEmbedError(text, title string) error {
	return c.RespondTextEmbed(text, &discordgo.MessageEmbed{Title: title, Color: EmbedColorError})
}
