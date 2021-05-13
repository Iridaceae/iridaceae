package pkg

import (
	"os"
	"os/exec"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// TODO: a better way to get general token.
func MakeTestSession() *discordgo.Session {
	botToken := os.Getenv("CONCERTINA_AUTHTOKEN")
	// ensure sessions are established
	dg, err := discordgo.New("Bot " + botToken)
	if err != nil {
		panic(err)
	}
	return dg
}

func GetEnvOrDefault(env, def string) string {
	v := os.Getenv(env)
	if v == "" && def != "" {
		v = def
	}
	return v
}

func GetRevision() string {
	// check for errors instead of printing to os.Stdout
	stdout, _ := exec.Command("git", "rev-parse", "--short", "HEAD").Output()
	return strings.ReplaceAll(string(stdout), "\n", "")
}

func GetVersion() string {
	return "v1"
}
