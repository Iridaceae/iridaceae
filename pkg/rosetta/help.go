// Package rosetta is iridaceae internal command parser. Greatly inspired by
// luis/dgc and zekroTJA/sireikan with some modification.
package rosetta

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Iridaceae/iridaceae/internal/helpers"

	"github.com/bwmarrin/discordgo"
)

type DefaultHelpCommand struct{}

func (d *DefaultHelpCommand) GetInvokers() []string {
	return []string{"help", "h", "?", "man"}
}

func (d *DefaultHelpCommand) GetDescription() string {
	return "display list of commands or help of specific command"
}

func (d *DefaultHelpCommand) GetUsage() string {
	return "`help` - display command list\n" + "`help <command>` - display help of a specific command"
}

func (d *DefaultHelpCommand) GetGroup() string {
	return GroupGeneral
}

func (d *DefaultHelpCommand) GetDomain() string {
	return "rs.etc.help"
}

func (d *DefaultHelpCommand) GetSubPermissionRules() []SubPermission {
	return nil
}

func (d *DefaultHelpCommand) IsExecutableInDM() bool {
	return true
}

func (d *DefaultHelpCommand) Exec(ctx Context) error {
	embed := &discordgo.MessageEmbed{
		Color:     helpers.EmbedColorDefault,
		Fields:    make([]*discordgo.MessageEmbedField, 0),
		Timestamp: time.RFC3339,
		Footer:    &discordgo.MessageEmbedFooter{Text: "with :hearts: and :coffee: by iridaceae"},
	}

	rr, _ := ctx.GetObject(ObjectMapKeyRouter).(Router)

	if ctx.GetArguments().Len() == 0 {
		// we will send a general help to dms.
		cmds := make(map[string][]Command)
		for _, c := range rr.GetCommandInstances() {
			group := c.GetGroup()
			if _, ok := cmds[group]; !ok {
				cmds[group] = make([]Command, 0)
			}
			cmds[group] = append(cmds[group], c)
		}

		embed.Title = fmt.Sprintf("Command List for %s", ctx.GetGuild().Name)

		for group, groupCmds := range cmds {
			cmdInfo := ""
			for i, c := range groupCmds {
				cmdInfo += fmt.Sprintf("%d. `%s` - *%s* `[%s]`", i, c.GetInvokers()[0], c.GetDescription(), c.GetDomain())
			}
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{Name: group, Value: cmdInfo})
		}
	} else {
		// specific commands we want to render
		invoke := ctx.GetArguments().Get(0).String()
		cmd, ok := rr.GetCommand(invoke)
		if !ok {
			_, err := ctx.RespondEmbedError(fmt.Sprintf("No command was found with given invoke `%s`.", invoke), ErrInvokeDoesNotExists)
			return err
		}

		embed.Title = "Command Description"
		description := cmd.GetDescription()
		if description == "" {
			description = "`no description`"
		}

		usage := cmd.GetUsage()
		if usage == "" {
			usage = "`no usage information`"
		}

		embed.Fields = []*discordgo.MessageEmbedField{
			{
				Name:   "Invokers",
				Value:  strings.Join(cmd.GetInvokers(), " "),
				Inline: true,
			},
			{
				Name:   "Group",
				Value:  cmd.GetGroup(),
				Inline: true,
			},
			{
				Name:   "Domain",
				Value:  cmd.GetDomain(),
				Inline: true,
			},
			{
				Name:   "IsDM-able",
				Value:  strconv.FormatBool(cmd.IsExecutableInDM()),
				Inline: true,
			},
			{
				Name:  "Description",
				Value: description,
			},
			{
				Name:  "Usage",
				Value: usage,
			},
		}

		if spr := cmd.GetSubPermissionRules(); spr != nil {
			txt := "*`[E]` means explicit permissions, or permission must be explicitly allowed and cannot be wild-carded.\n" +
				"`[NE]` means non-explicit permissions, or wildcards will apply to this sub permission.*\n\n"

			for _, rule := range spr {
				expl := "NE"
				if rule.Explicit {
					expl = "E"
				}
				txt = fmt.Sprintf("%s`[%s]` %s - *%s*\n", txt, expl, getTermAssembly(cmd, rule.Term), rule.Description)
			}

			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:  "Sub Permission Rules",
				Value: txt,
			})
		}
	}

	channel, err := ctx.GetSession().UserChannelCreate(ctx.GetUser().ID)
	if err != nil {
		return err
	}
	_, err = ctx.GetSession().ChannelMessageSendEmbed(channel.ID, embed)
	if err != nil {
		if strings.Contains(err.Error(), `{"code": 50007, "message": "Cannot send messages to this user"}`) {
			embed.Footer = &discordgo.MessageEmbedFooter{
				Text: "This message appears in DMs, but you have disabled `receiving DMs from server members`",
			}
			_, err = ctx.GetSession().ChannelMessageSendEmbed(ctx.GetChannel().ID, embed)
		}
	}
	return err
}

// parse our SubPermissions rules to our domain.
func getTermAssembly(cmd Command, term string) string {
	if strings.HasPrefix(term, "/") {
		return term[1:]
	}
	return cmd.GetDomain() + "." + term
}
