package ratelimit

import (
	"github.com/Iridaceae/iridaceae/pkg/rosetta"

	"github.com/bwmarrin/discordgo"
)

type TestContext struct {
	chanType discordgo.ChannelType
	gid      string
	uid      string
}

func (tc *TestContext) GetObject(key string) (value interface{}) { return nil }

func (tc *TestContext) SetObject(key string, value interface{}) {}

func (tc *TestContext) GetSession() *discordgo.Session {
	return nil
}

func (tc *TestContext) GetArguments() *rosetta.Arguments {
	return nil
}

func (tc *TestContext) GetChannel() *discordgo.Channel {
	return &discordgo.Channel{Type: tc.chanType}
}

func (tc *TestContext) GetMessage() *discordgo.Message {
	return nil
}

func (tc *TestContext) GetGuild() *discordgo.Guild {
	return &discordgo.Guild{ID: tc.gid}
}

func (tc *TestContext) GetUser() *discordgo.User {
	return &discordgo.User{ID: tc.uid}
}

func (tc *TestContext) GetMember() *discordgo.Member {
	return nil
}

func (tc *TestContext) IsDM() bool {
	return false
}

func (tc *TestContext) IsEdit() bool {
	return false
}

func (tc *TestContext) RespondText(content string) (*discordgo.Message, error) {
	return nil, nil
}

func (tc *TestContext) RespondEmbed(embed *discordgo.MessageEmbed) (*discordgo.Message, error) {
	return nil, nil
}

func (tc *TestContext) RespondEmbedError(title string, err error) (*discordgo.Message, error) {
	return nil, nil
}

func (tc *TestContext) RespondTextEmbed(content string, embed *discordgo.MessageEmbed) (*discordgo.Message, error) {
	return nil, nil
}

func (tc *TestContext) RespondTextEmbedError(title, content string, err error) (*discordgo.Message, error) {
	return nil, nil
}

type TestCmdNotImplemented struct{}

func (t *TestCmdNotImplemented) GetInvokers() []string {
	return nil
}

func (t *TestCmdNotImplemented) GetDescription() string {
	return ""
}

func (t *TestCmdNotImplemented) GetUsage() string {
	return ""
}

func (t *TestCmdNotImplemented) GetGroup() string {
	return ""
}

func (t *TestCmdNotImplemented) GetDomain() string {
	return ""
}

func (t *TestCmdNotImplemented) GetSubPermissionRules() []rosetta.SubPermission {
	return nil
}

func (t *TestCmdNotImplemented) IsExecutableInDM() bool {
	return false
}

func (t *TestCmdNotImplemented) Exec(ctx rosetta.Context) error {
	return rosetta.ErrCommandExec
}
