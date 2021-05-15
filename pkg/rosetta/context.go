package rosetta

import (
	"fmt"
	"sync"

	"github.com/bwmarrin/discordgo"
)

// Context is an interface representing information about a message and environment
// where this message was created and passed to middleware and command router.
type Context interface {
	ObjectMap

	// GetSession returns our instance of discordgo.Session.
	GetSession() *discordgo.Session

	// GetArguments returns our Arguments list and parsed Command Arguments.
	GetArguments() *Arguments

	// GetChannel returns the channel where message is sent.
	GetChannel() *discordgo.Channel

	// GetMessage returns the content of sent message.
	GetMessage() *discordgo.Message

	// GetGuild returns guild objects where command was sent.
	// We can use this later for logging purposes, update databases, etc.
	GetGuild() *discordgo.Guild

	// GetUser returns said user who invokes the command.
	GetUser() *discordgo.User

	// GetMember returns the member object of the author of the message.
	GetMember() *discordgo.Member

	// IsDM returns true if contextImpl is sent in a dms or group dms, false otherwise
	IsDM() bool

	// IsEdit returns true if event is a *discordgo.MessageUpdate event.
	IsEdit() bool

	// RespondText wraps around responses of given text message.
	RespondText(content string) (*discordgo.Message, error)

	// RespondEmbed responds with the given embed message.
	RespondEmbed(embed *discordgo.MessageEmbed) (*discordgo.Message, error)

	// RespondEmbedError responds with the given error in a embed message.
	RespondEmbedError(title string, err error) (*discordgo.Message, error)
}

// contextImpl is our default implementation of Context.
type contextImpl struct {
	isDM      bool
	isEdit    bool
	router    Router
	args      *Arguments
	objectMap *sync.Map
	session   *discordgo.Session
	message   *discordgo.Message
	guild     *discordgo.Guild
	channel   *discordgo.Channel
	member    *discordgo.Member
}

func (c *contextImpl) GetObject(key string) (value interface{}) {
	var ok bool
	if c.objectMap != nil {
		value, ok = c.objectMap.Load(key)
	}
	// if our internal object map doesn't contain the key, then get from di.Container.
	if !ok {
		value = c.router.GetObject(key)
	}
	return
}

func (c *contextImpl) SetObject(key string, value interface{}) {
	c.objectMap.Store(key, value)
}

func (c *contextImpl) GetSession() *discordgo.Session {
	return c.session
}

func (c *contextImpl) GetArguments() *Arguments {
	return c.args
}

func (c *contextImpl) GetChannel() *discordgo.Channel {
	return c.channel
}

func (c *contextImpl) GetMessage() *discordgo.Message {
	return c.message
}

func (c *contextImpl) GetGuild() *discordgo.Guild {
	return c.guild
}

func (c *contextImpl) GetUser() *discordgo.User {
	return c.message.Author
}

func (c *contextImpl) GetMember() *discordgo.Member {
	return c.member
}

func (c *contextImpl) IsDM() bool {
	return c.isDM
}

func (c *contextImpl) IsEdit() bool {
	return c.isEdit
}

func (c *contextImpl) RespondText(content string) (*discordgo.Message, error) {
	return c.session.ChannelMessageSend(c.channel.ID, content)
}

func (c *contextImpl) RespondEmbed(embed *discordgo.MessageEmbed) (*discordgo.Message, error) {
	return c.session.ChannelMessageSendEmbed(c.channel.ID, embed)
}

func (c *contextImpl) RespondEmbedError(title string, err error) (*discordgo.Message, error) {
	return c.session.ChannelMessageSendEmbed(c.channel.ID, &discordgo.MessageEmbed{Title: title, Description: fmt.Sprintf("*%s*", err.Error()), Color: EmbedColorError})
}
