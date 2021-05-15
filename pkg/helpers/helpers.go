package helpers

import (
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// DomainRegex follows our same rule for configparser.
const DomainRegex string = "^(([\\w\\.])+(\\.)([\\w]){2,4}([\\w]*))*$"

func MakeTestSession() *discordgo.Session {
	// TODO: a better way to get general token.
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

func CheckDomainCompile(name string) bool {
	if _, err := regexp.MatchString(DomainRegex, name); err != nil {
		return false
	}
	return true
}
