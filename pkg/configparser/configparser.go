// Package configparser defines some default configs handler including a configparser parser with ability to update value dynamically
package configparser

import (
	"fmt"
)

const OptionsRegex string = "^(([\\w\\.])+(\\.)([\\w]){2,4}([\\w]*))*$"

var (
	ErrEmptyValue          = fmt.Errorf("empty strings")
	ErrInvalidFormat       = fmt.Errorf("invalid format")
	ErrInvalidOptionsMatch = fmt.Errorf("invalid options match")
)

// Source acts as a generic type for different source of configs.
// ex: env, yaml, toml. refers to EnvSource for envars parsing.
// Source also have the ability to marshall and unmarshal given config to file.
type Source interface {

	// GetValue will return our requested value from key.
	// This will throw error if none was found or given key doesn't
	// follow our regex parsing.
	GetValue(key string) (interface{}, error)

	// Name returns name of given source.
	Name() string

	// Marshal serializes a data stream to an io.Writer from
}

// ConfigManager holds types for generic managers to generate configs.
type ConfigManager struct {
	sources []Source
	Options map[string]*Options
}

// NewConfigManager makes a configs manager.
func NewConfigManager() *ConfigManager {
	return &ConfigManager{Options: make(map[string]*Options)}
}

// AddSource allows users to append given configparser source to the manager.
func (c *ConfigManager) AddSource(source Source) {
	c.sources = append(c.sources, source)
}

// Register will add given configs to the general manager.
func (c *ConfigManager) Register(name, desc string, defaultValue interface{}) (*Options, error) {
	if _, err := matchOptionsRegex(name); err != nil {
		return nil, ErrInvalidFormat
	}
	opt := &Options{
		Name:         name,
		Description:  desc,
		DefaultValue: defaultValue,
		Manager:      c,
	}
	c.Options[name] = opt
	return opt, nil
}

// Load handles configs func LoadValue directly.
func (c *ConfigManager) Load() {
	for _, v := range c.Options {
		v.LoadValue()
	}
}

func (c *ConfigManager) Clear() {
	c.sources = make([]Source, 0)
	c.Options = make(map[string]*Options)
}
