package jog

import (
	"os"

	"github.com/bwmarrin/discordgo"
)

var (
	dg       *discordgo.Session // stores a global discordgo user session.
	dgBot    *discordgo.Session // stores a global discordgo bot session.
	envToken = os.Getenv("")
)
