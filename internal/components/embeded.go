package components

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

// EmbedMessage wraps discordgo.Message with some more features.
type EmbedMessage struct {
	*discordgo.Message
	s   *discordgo.Session
	err error
}

// DeleteAfter deletes message after a specified duration when the message
// still exists and returned sets error to EmbedMessage.
func (e *EmbedMessage) DeleteAfter(d time.Duration) *EmbedMessage {
	if e.Message != nil {
		time.AfterFunc(d, func() {
			e.err = e.s.ChannelMessageDelete(e.ChannelID, e.ID)
		})
	}
	return e
}

// Error returns embed error.
func (e *EmbedMessage) Error() error {
	return e.err
}

// Edit updates given embed message with given content replace internal message
// and error of the embed instance.
func (e *EmbedMessage) Edit(content, title string, color int) *EmbedMessage {
	newEmbed := &discordgo.MessageEmbed{
		Type:        "rich",
		Title:       title,
		Description: content,
		Color:       color,
		Timestamp:   time.Now().Format(time.RFC3339),
		Footer:      &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("edited by %s", e.Author.Username)},
	}
	if color == 0 {
		newEmbed.Color = EmbedColorDefault
	}
	return e.EditRaw(newEmbed)
}

// EditRaw updates current embed message with given raw embed then replace
// internal message and error of this embed instance.
func (e *EmbedMessage) EditRaw(embed *discordgo.MessageEmbed) *EmbedMessage {
	e.Message, e.err = e.s.ChannelMessageEditEmbed(e.ChannelID, e.ID, embed)
	return e
}

// SendEmbed creates a discordgo.MessageEmbed afrom passed content, title, and color
// then send it to specific channel.
// If no color is specified, then use EmbedColorViolet.
func SendEmbed(s *discordgo.Session, channelID, content, title string, color int) *EmbedMessage {
	e := &discordgo.MessageEmbed{
		Type:        "rich",
		Description: content,
		Color:       color,
		Timestamp:   time.Now().Format(time.RFC3339),
		Title:       title,
	}
	if color == 0 {
		e.Color = EmbedColorViolet
	}
	return SendEmbedRaw(s, channelID, e)
}

// SendEmbedRaw passed embed to a channel and set occurred error to internal errors.
func SendEmbedRaw(s *discordgo.Session, channelID string, embed *discordgo.MessageEmbed) *EmbedMessage {
	msg, err := s.ChannelMessageSendEmbed(channelID, embed)
	return &EmbedMessage{msg, s, err}
}

// SendEmbedError will send given error to user about a specific command errors.
func SendEmbedError(s *discordgo.Session, channelID string, err error) *EmbedMessage {
	e := &discordgo.MessageEmbed{
		Type:        "rich",
		Title:       "Error",
		Description: err.Error(),
		Timestamp:   time.Now().Format(time.RFC3339),
		Color:       EmbedColorError,
		Footer:      &discordgo.MessageEmbedFooter{Text: "Please contact @aarnphm for help"},
	}
	return SendEmbedRaw(s, channelID, e)
}

// SendEmbedComplex sends a embed message with user mentions that allows to notify users about certain tasks.
func (e *EmbedMessage) SendEmbedComplex(s *discordgo.Session, channelID, content, title string) *EmbedMessage {
	embed := &discordgo.MessageEmbed{
		Type:        "rich",
		Title:       title,
		Description: content,
		Timestamp:   time.Now().Format(time.RFC3339),
		Color:       EmbedColorDefault,
	}
	mentioned := fmt.Sprintf("%s\n", e.Author.Mention())
	return SendEmbedComplexRaw(s, embed, channelID, mentioned)
}

// SendEmbedComplexRaw takes given mentions and embed then streamline to given channel, returns correspondingly
// EmbedMessage instance with internal error.
func SendEmbedComplexRaw(s *discordgo.Session, embed *discordgo.MessageEmbed, channelID, mention string) *EmbedMessage {
	msg, err := s.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{Content: mention, Embed: embed})
	return &EmbedMessage{msg, s, err}
}
