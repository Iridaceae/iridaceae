// Package util provides a general utilities functions that can be used globally.
package util

import (
	"os"
	"os/exec"
	"strings"
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
