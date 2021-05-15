// Package acceptmsg provides a message models for
// discordgo which can be either accepted or declined via reactions.
package acceptmsg

import (
	"github.com/bwmarrin/discordgo"

	helpers2 "github.com/Iridaceae/iridaceae/internal/helpers"
)

const (
	acceptMessageEmoteAccept  = "✅"
	acceptMessageEmoteDecline = "❌"
)

type ActionHandler func(*discordgo.Message)

// AcceptMessage extends discordgo.Message to build and send an AcceptMessage.
type AcceptMessage struct {
	*discordgo.Message
	Session        *discordgo.Session
	Embed          *discordgo.MessageEmbed
	UserID         string
	DeleteMsgAfter bool
	AcceptFunc     ActionHandler
	DeclineFunc    ActionHandler
	eventListener  func()
}

// New creates an empty instance of AcceptMessage.
func New() *AcceptMessage {
	return new(AcceptMessage)
}

// WithSession set a discordgo.Session.
func (a *AcceptMessage) WithSession(s *discordgo.Session) *AcceptMessage {
	a.Session = s
	return a
}

// WithEmbed sets the embed instance.
func (a *AcceptMessage) WithEmbed(e *discordgo.MessageEmbed) *AcceptMessage {
	a.Embed = e
	return a
}

// WithContent creates an embed with default color and specified content as descriptions.
func (a *AcceptMessage) WithContent(content string) *AcceptMessage {
	a.Embed = &discordgo.MessageEmbed{
		Color:       helpers2.EmbedColorDefault,
		Description: content,
	}
	return a
}

// AcceptOnlyUser specifies only determined users can have inputs.
func (a *AcceptMessage) AcceptOnlyUser(userID string) *AcceptMessage {
	a.UserID = userID
	return a
}

// DeleteAfterAnswer enables embed message to be delete after users' answer.
func (a *AcceptMessage) DeleteAfterAnswer() *AcceptMessage {
	a.DeleteMsgAfter = true
	return a
}

// OnAccept specifies action handler to be executed if accept.
func (a *AcceptMessage) OnAccept(onAcc ActionHandler) *AcceptMessage {
	a.AcceptFunc = onAcc
	return a
}

// OnDecline specifies action handler to be executed if decline.
func (a *AcceptMessage) OnDecline(onDec ActionHandler) *AcceptMessage {
	a.DeclineFunc = onDec
	return a
}

// Send pushes accept message into a channel and setup listener handler for reactions.
func (a *AcceptMessage) Send(channelID string) (*AcceptMessage, error) {
	if a.Session == nil {
		return nil, helpers2.ErrSessionNotDefined
	}
	if a.Embed == nil {
		return nil, helpers2.ErrEmbedNotDefined
	}

	msg, _ := a.Session.ChannelMessageSendEmbed(channelID, a.Embed)

	a.Message = msg
	_ = a.Session.MessageReactionAdd(channelID, msg.ID, acceptMessageEmoteAccept)
	_ = a.Session.MessageReactionAdd(channelID, msg.ID, acceptMessageEmoteDecline)

	a.eventListener = a.Session.AddHandler(func(s *discordgo.Session, e *discordgo.MessageReactionAdd) {
		if e.MessageID != msg.ID {
			return
		}

		if e.UserID != a.Session.State.User.ID {
			_ = a.Session.MessageReactionRemove(a.ChannelID, a.ID, e.Emoji.Name, e.UserID)
		}
		if e.UserID == s.State.User.ID || (a.UserID != "" && a.UserID != e.UserID) {
			return
		}

		if e.Emoji.Name != acceptMessageEmoteAccept && e.Emoji.Name != acceptMessageEmoteDecline {
			return
		}

		switch e.Emoji.Name {
		case acceptMessageEmoteDecline:
			if a.DeclineFunc != nil {
				a.DeclineFunc(msg)
			}
		case acceptMessageEmoteAccept:
			if a.AcceptFunc != nil {
				a.AcceptFunc(msg)
			}
		}

		a.eventListener()
		if a.DeleteMsgAfter {
			_ = a.Session.ChannelMessageDelete(channelID, msg.ID)
		} else {
			_ = a.Session.MessageReactionsRemoveAll(channelID, msg.ID)
		}
	})
	return a, nil
}
