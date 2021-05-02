// Package configparser defines some default configs handler including a configparser parser with ability to update value dynamically
package configparser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const optionsFmtRegex string = "^(([\\w\\.])+(\\.)([\\w]){2,4}([\\w]*))*$"

var (
	ErrEmptyValue          = fmt.Errorf("empty strings")
	ErrInvalidFormat       = fmt.Errorf("invalid format")
	ErrInvalidOptionsMatch = fmt.Errorf("invalid options match")
)

// Source acts as a generic type for different source of configs.
// ex: env, yaml, toml. refers to EnvSource for envars parsing.
type Source interface {
	GetValue(key string) (interface{}, error)
	Name() string
}

type Options struct {
	// Name will have format iris.option1.option2
	Name         string
	Description  string
	DefaultValue interface{}
	LoadedValue  interface{}
	Manager      *Manager

	ConfigSource Source
}

// Manager holds types for generic managers to generate configs.
type Manager struct {
	sources []Source
	Options map[string]*Options
}

// NewManager makes a configs manager.
func NewManager() *Manager {
	return &Manager{Options: make(map[string]*Options)}
}

// AddSource allows users to append given configparser source to the manager.
func (c *Manager) AddSource(source Source) {
	c.sources = append(c.sources, source)
}

// Register will add given configs to the general manager.
func (c *Manager) Register(name, desc string, defaultValue interface{}) (*Options, error) {
	_, err := matchOptionsRegex(name)
	if err != nil {
		return &Options{}, err
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
func (c *Manager) Load() {
	for _, v := range c.Options {
		v.LoadValue()
	}
}

// LoadValue will load given values if exists, otherwise use default ones.
func (opt *Options) LoadValue() {
	def := opt.DefaultValue
	opt.ConfigSource = nil

	for i := len(opt.Manager.sources) - 1; i >= 0; i-- {
		source := opt.Manager.sources[i]
		// v would be value from given source, check envsource.go for examples
		v, _ := source.GetValue(opt.Name)

		if v != nil {
			def = v
			opt.ConfigSource = source
			break
		}
	}

	if opt.DefaultValue != nil {
		if _, ok := opt.DefaultValue.(int); ok {
			def = interface{}(toIntVal(def))
		} else if _, ok := opt.DefaultValue.(bool); ok {
			def = interface{}(toBoolVal(def))
		}
	}

	opt.LoadedValue = def
}

// UpdateValue updates loaded value.
func (opt *Options) UpdateValue(val interface{}) {
	switch val.(type) {
	case bool:
		opt.LoadedValue = toBoolVal(val)
	case string:
		opt.LoadedValue = toStrVal(val)
	case int:
		opt.LoadedValue = toIntVal(val)
	case float64:
		opt.LoadedValue = toFloat64Val(val)
	}
}

// GetString are a getter string for &Options.LoadedValue.
func (opt *Options) GetString() string {
	return toStrVal(opt.LoadedValue)
}

// GetInt are a getter int for &Options.LoadedValue.
func (opt *Options) GetInt() int {
	return toIntVal(opt.LoadedValue)
}

// GetBool are a getter bool for &Options.LoadedValue.
func (opt *Options) GetBool() bool {
	return toBoolVal(opt.LoadedValue)
}

// GetFloat are a getter float64 for &Options.LoadedValue.
func (opt *Options) GetFloat() float64 {
	return toFloat64Val(opt.LoadedValue)
}

func toStrVal(i interface{}) string {
	switch t := i.(type) {
	case string:
		return t
	case int:
		return strconv.FormatInt(int64(t), 10)
	case fmt.Stringer:
		return t.String()
	}
	return ""
}

func toIntVal(i interface{}) int {
	switch t := i.(type) {
	case string:
		n, _ := strconv.ParseInt(t, 10, 64)
		return int(n)
	case int:
		return t
	}
	return 0
}

func toFloat64Val(i interface{}) float64 {
	switch t := i.(type) {
	case string:
		n, _ := strconv.ParseFloat(t, 64)
		return n
	case int:
		return float64(t)
	case float64:
		return t
	}
	return 0
}

func toBoolVal(i interface{}) bool {
	switch t := i.(type) {
	case string:
		lower := strings.ToLower(strings.TrimSpace(t))
		// NOTE: regex match
		if lower == "true" || lower == "yes" || lower == "on" || lower == "enabled" || lower == "1" {
			return true
		}
		return false
	case int:
		return t > 0
	case bool:
		return t
	}

	return false
}

func matchOptionsRegex(key string) (bool, error) {
	b, _ := regexp.MatchString(optionsFmtRegex, key)
	if b {
		return b, nil
	}
	return b, ErrInvalidOptionsMatch
}
