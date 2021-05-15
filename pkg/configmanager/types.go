package configmanager

import (
	"fmt"
	"io"
)

const OptionsRegex string = "^(([\\w\\.])+(\\.)([\\w]){2,4}([\\w]*))*$"

var (
	ErrEmptyValue          = fmt.Errorf("empty strings")
	ErrInvalidFormat       = fmt.Errorf("invalid format")
	ErrInvalidOptionsMatch = fmt.Errorf("invalid options match")
)

// Source acts as a generic type for different source of configs.
type Source interface {

	// GetValue will return our requested value from key.
	// This will throw error if none was found or given key doesn't
	// follow our regex parsing.
	GetValue(key string) (interface{}, error)

	// Name returns name of given source.
	Name() string
}

// Parser allows us to handle read and write to and from file.
type Parser interface {
	// Marshal serialize data stream into io.Writer from our config
	// instance and return errors during serialization.
	Marshal(w io.Writer, c *Config) error

	// Unmarshal deserializes a config instance from io.Reader and
	// return our config instance and error occurred during serialization.
	Unmarshal(r io.Reader) (*Config, error)
}

type Manager interface {

	// AddSource allows users to append given configparser source to the manager.
	AddSource(s Source)

	// Register will add given configs to the general manager.
	Register(name, desc string, defaultValue interface{}) (*Options, error)

	// Load will configure our options value into given manager.
	Load()

	// Reset will reset our sources and options mapping.
	Reset()
}

// Config is specific configuration for iridaceae.
type Config struct {
	// TODO: generate these from files for extension.
	Version    int         `yaml:"version" json:"version"`
	Discord    *Discord    `yaml:"discord" json:"discord"`
	Permission *Permission `yaml:"permission" json:"permission"`
	Logging    *Logging    `yaml:"logging" json:"logging"`
	Metrics    *Metrics    `yaml:"metrics" json:"metrics"`
}

type Discord struct {
	GlobalRateLimit *GlobalRateLimit ` yaml:"globalratelimit" json:"globalratelimit"`
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
