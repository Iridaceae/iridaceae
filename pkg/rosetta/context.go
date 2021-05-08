package rosetta

import (
	"github.com/bwmarrin/discordgo"
)

// ExecutionHandler represents a handler for a context execution.
type ExecutionHandler func(ctx *Context, meta ...interface{})

// ResponseModal represents response from rosetta to given handler.
type ResponseModal interface {

	// RespondText wraps around responses of given text message.
	RespondText(text string) error

	// RespondEmbed responds with the given embed message.
	RespondEmbed(embed *discordgo.MessageEmbed) error

	// RespondEmbedError responds with the given error in a embed message.
	RespondEmbedError(err error) error

	// RespondTextEmbed responds with given text and embed message.
	RespondTextEmbed(text string, embed *discordgo.MessageEmbed) error

	// RespondTextEmbedError responds given error to users with embed message.
	RespondTextEmbedError(text, title string, err error) error
}

// Context represents the context for a command event.
type Context struct {
	Session    *discordgo.Session
	Event      *discordgo.MessageCreate
	Channel    *discordgo.Channel
	Arguments  *Arguments
	Router     *Router
	Command    *Command
	ObjectsMap *ObjectsMap
}

func (c *Context) RespondText(text string) error {
	_, err := c.Session.ChannelMessageSend(c.Event.ChannelID, text)
	return err
}

func (c *Context) RespondEmbed(embed *discordgo.MessageEmbed) error {
	_, err := c.Session.ChannelMessageSendEmbed(c.Event.ChannelID, embed)
	return err
}

func (c *Context) RespondEmbedError(e error) error {
	_, err := c.Session.ChannelMessageSendEmbed(c.Event.ChannelID, &discordgo.MessageEmbed{Title: "Error", Description: e.Error(), Color: EmbedColorError})
	return err
}

func (c *Context) RespondTextEmbed(text string, embed *discordgo.MessageEmbed) error {
	_, err := c.Session.ChannelMessageSendComplex(c.Event.ChannelID, &discordgo.MessageSend{
		Content: text,
		Embed:   embed,
	})
	return err
}

func (c *Context) RespondTextEmbedError(text, title string, err error) error {
	return c.RespondTextEmbed(text, &discordgo.MessageEmbed{Title: title, Description: err.Error(), Color: EmbedColorError})
}
