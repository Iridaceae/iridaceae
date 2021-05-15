package configparser

import (
	"io"

	"github.com/Iridaceae/iridaceae/internal/helpers"
)

// Parser allows us to handle read and write to and from file.
type Parser interface {
	// Marshal serialize data stream into io.Writer from our configparser
	// instance and return errors during serialization.
	Marshal(w io.Writer, c Settings) error

	// Unmarshal deserializes a configparser instance from io.Reader and
	// return our configparser instance and error occurred during serialization.
	Unmarshal(r io.Reader) (Settings, error)
}

// Settings is specific configuration for iridaceae.
type Settings struct {
	// TODO: generate these from files for extension.
	Version    int         `yaml:"version" json:"version"`
	Discord    *Discord    `yaml:"discord" json:"discord"`
	Permission *Permission `yaml:"permission" json:"permission"`
	Logging    *Logging    `yaml:"logging" json:"logging"`
	Metrics    *Metrics    `yaml:"metrics" json:"metrics"`
}

type Discord struct {
	Prefix          string           `yaml:"prefix" json:"prefix"`
	GlobalRateLimit *GlobalRateLimit `yaml:"globalratelimit" json:"globalratelimit"`
}

type GlobalRateLimit struct {
	Burst       int `json:"burst"`
	Restoration int `json:"restoration"`
}

type Permission struct {
	User  []string `yaml:"user" json:"user"`
	Admin []string `yaml:"admin" json:"admin"`
}

type Logging struct {
	Enabled bool `yaml:"enabled" json:"enabled"`
	Level   int  `yaml:"level" json:"level"`
}

type Metrics struct {
	Enabled bool   `yaml:"enabled" json:"enabled"`
	Address string `yaml:"address" json:"address"`
}

func GetSettings() *Settings {
	return &Settings{
		Version:    1,
		Discord:    &Discord{Prefix: "ir!", GlobalRateLimit: &GlobalRateLimit{Burst: 2, Restoration: 10}},
		Permission: &Permission{User: helpers.DefaultUserRules, Admin: helpers.DefaultAdminRules},
		Logging:    &Logging{Enabled: true, Level: 1},
		Metrics:    &Metrics{Enabled: true, Address: ":9000"},
	}
}
