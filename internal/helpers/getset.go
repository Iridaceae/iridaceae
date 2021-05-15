package helpers

import (
	"os"
	"os/exec"
	"regexp"
	"strings"
)

const (
	// DomainRegex follows our same rule for configmanager.
	DomainRegex string = "^(([\\w\\.])+(\\.)([\\w]){2,4}([\\w]*))*$"
	// BaseAuthURLTemplate is our default API invitation link.
	BaseAuthURLTemplate string = "https://discord.com/api/oauth2/authorize?client_id=%s&scope=bot&permissions=%d"
)

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
