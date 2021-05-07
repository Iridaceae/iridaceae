package jog

import (
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// RegisterDefaultHelpCommand registers the default help command.
func (r *Router) RegisterDefaultHelpCommand(session *discordgo.Session, limiter *RateLimiter) {
	r.InitializeStorage("jog_helpMessages")

	// Init reaction add listener.
	session.AddHandler(func(session *discordgo.Session, event *discordgo.MessageReactionAdd) {
		channelID := event.ChannelID
		messageID := event.MessageID
		userID := event.UserID

		// Check if the reaction is added by the bot.
		if event.UserID == session.State.User.ID {
			return
		}

		// check if the message is a help message.
		raw, ok := r.Storage["jog_helpMessages"].Get(channelID + ":" + messageID + ":" + event.UserID)
		if !ok {
			return
		}

		page, ok := raw.(int)
		if !ok {
			return
		}
		if page <= 0 {
			return
		}

		// check which reaction was added.
		reactionName := event.Emoji.Name
		switch reactionName {
		case "⬅️":
			// Update help message
			embed, newPage := renderDefaultGeneralHelpEmbed(r, page-1)
			page = newPage
			if _, err := session.ChannelMessageEditEmbed(channelID, messageID, embed); err != nil {
				return
			}

			// Remove reaction.
			if err := session.MessageReactionRemove(channelID, messageID, reactionName, userID); err != nil {
				return
			}
		case "❌":
			// Delete help message
			if err := session.ChannelMessageDelete(channelID, messageID); err != nil {
				return
			}
		case "➡️":
			// Update help message
			embed, newPage := renderDefaultGeneralHelpEmbed(r, page+1)
			page = newPage
			if _, err := session.ChannelMessageEditEmbed(channelID, messageID, embed); err != nil {
				return
			}

			// Remove reaction.
			if err := session.MessageReactionRemove(channelID, messageID, reactionName, userID); err != nil {
				return
			}
		}
		// Update stores page
		r.Storage["jog_helpMessages"].Set(channelID+":"+messageID+":"+event.UserID, page)
	})

	// Register default help commands.
	r.RegisterCmd(&Command{
		Name:        "help",
		Aliases:     []string{"h", "?", "help"},
		Description: "Lists all available commands or displays information about a specific command",
		Usage:       "help/h/? [command_name]",
		Example:     "help pom",
		IgnoreCase:  true,
		RateLimiter: limiter,
		Handler:     generalHelpCommand,
	})
}

func generalHelpCommand(ctx *Context) {
	if ctx.Arguments.Len() > 0 {
		specificHelpCommand(ctx)
		return
	}

	// some variables here.
	channelID := ctx.Event.ChannelID
	session := ctx.Session

	// send general help embed
	embed, _ := renderDefaultGeneralHelpEmbed(ctx.Router, 1)
	message, _ := ctx.Session.ChannelMessageSendEmbed(channelID, embed)

	// Add reactions to helper.
	if err := session.MessageReactionAdd(channelID, message.ID, "⬅️"); err != nil {
		return
	}
	if err := session.MessageReactionAdd(channelID, message.ID, "❌"); err != nil {
		return
	}
	if err := session.MessageReactionAdd(channelID, message.ID, "➡️"); err != nil {
		return
	}
	ctx.Router.Storage["jog_helpMessages"].Set(channelID+":"+message.ID+":"+ctx.Event.Author.ID, 1)
}

func renderDefaultGeneralHelpEmbed(r *Router, page int) (*discordgo.MessageEmbed, int) {
	commands := r.Commands
	prefix := r.Prefixes[0]

	// get # of pages.
	pageAmount := int(math.Ceil(float64(len(commands)) / 5))
	if page > pageAmount {
		page = pageAmount
	}
	if page <= 0 {
		page = 1
	}

	// get slice of commands to display on a page.
	startIdx := (page - 1) * 5
	endIdx := startIdx + 5
	if page == pageAmount {
		endIdx = len(commands)
	}
	displayCmds := commands[startIdx:endIdx]

	// Prep the fields for the embed
	fields := make([]*discordgo.MessageEmbedField, 0, len(displayCmds))
	for idx, cmd := range displayCmds {
		fields[idx] = &discordgo.MessageEmbedField{
			Name:   cmd.Name,
			Value:  "`" + cmd.Description + "`",
			Inline: false,
		}
	}

	// return the embed and new page
	return &discordgo.MessageEmbed{
		Type:        "rich",
		Title:       "Command List (Page " + strconv.Itoa(page) + "/" + strconv.Itoa(pageAmount) + ")",
		Description: "Helpers for available commands. Type `" + prefix + "help [command_name] to find out more about a specific command.",
		Timestamp:   time.Now().Format(time.RFC3339),
		Color:       EmbedColorDefault,
		Fields:      fields,
	}, page
}

func specificHelpCommand(ctx *Context) {
	// defines command name.
	cmdNames := strings.Split(ctx.Arguments.Raw(), " ")
	var cmd *Command
	for idx, cmdName := range cmdNames {
		if idx == 0 {
			cmd = ctx.Router.GetCmd(cmdName)
			continue
		}
		if cmd == nil {
			break
		}
		cmd = cmd.GetSubCmd(cmdName)
	}

	// Send the embed message
	if _, err := ctx.Session.ChannelMessageSendEmbed(ctx.Event.ChannelID, renderDefaultSpecificHelpEmbed(ctx, cmd)); err != nil {
		return
	}
}

func renderDefaultSpecificHelpEmbed(ctx *Context, command *Command) *discordgo.MessageEmbed {
	prefix := ctx.Router.Prefixes[0]

	// check if command is valid.
	if command == nil {
		return &discordgo.MessageEmbed{
			Type:      "rich",
			Title:     "Error",
			Timestamp: time.Now().Format(time.RFC3339),
			Color:     EmbedColorError,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Message",
					Value:  "```Given command (" + command.Name + ")doesn't exist. Type `" + prefix + "help` for a list of available commands.```",
					Inline: false,
				},
			},
		}
	}

	// Define subcommands string.
	subCmds := "No sub commands"
	if len(command.SubCommands) > 0 {
		subCmdNames := make([]string, 0, len(command.SubCommands))
		for idx, subCmd := range command.SubCommands {
			subCmdNames[idx] = subCmd.Name
		}
		subCmds = "`" + strings.Join(command.Aliases, "`, `") + "`"
	}

	// define aliases
	aliases := "No aliases"
	if len(command.Aliases) > 0 {
		aliases = "`" + strings.Join(command.Aliases, "`, `") + "`"
	}

	// return embed.
	return &discordgo.MessageEmbed{
		URL:         "",
		Type:        "rich",
		Title:       "Command Information",
		Description: "Displaying the information for the `" + command.Name + "` command.",
		Timestamp:   time.Now().Format(time.RFC3339),
		Color:       0xffff00,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Name",
				Value:  "`" + command.Name + "`",
				Inline: false,
			},
			{
				Name:   "Subcommands",
				Value:  subCmds,
				Inline: false,
			},
			{
				Name:   "Aliases",
				Value:  aliases,
				Inline: false,
			},
			{
				Name:   "Description",
				Value:  "```" + command.Description + "```",
				Inline: false,
			},
			{
				Name:   "Usage",
				Value:  "```" + prefix + command.Usage + "```",
				Inline: false,
			},
			{
				Name:   "Example",
				Value:  "```" + prefix + command.Example + "```",
				Inline: false,
			},
		},
	}
}
