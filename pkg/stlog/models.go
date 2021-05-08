package stlog

import "github.com/rs/zerolog"

type LogCtxKey string

var cfg LogConfig

// String is human readable representation of a context key.
func (l LogCtxKey) String() string {
	return "mw-" + string(l)
}

// Logger defines a default logger for iris that wraps around rs/zerolog.
type Logger struct {
	Level         int
	Version       string
	Revision      string
	StdLog        zerolog.Logger
	ErrLog        zerolog.Logger
	dynamicFields []interface{}
}

// LogConfig defines config for Logger.
type LogConfig struct {
	name         string
	level        int
	staticFields []interface{}
	configured   bool
}

func cfgSetup(name string, staticFields []interface{}) {
	if cfg.configured {
		return
	}

	cfg.name = name
	cfg.staticFields = append(cfg.staticFields, staticFields...)
	cfg.configured = true
}

// SetDynaFields acts as a receiver instance that will always append these key-value pairs to the output.
func (i *Logger) SetDynaFields(dynamicFields ...interface{}) {
	i.dynamicFields = make([]interface{}, 2)
	i.dynamicFields = append(i.dynamicFields, dynamicFields...)
}

// AddDynaFields acts as a receiver instance that add given key-value pairs to current logger.
func (i *Logger) AddDynaFields(key, value interface{}) {
	i.dynamicFields = append(i.dynamicFields, []interface{}{key, value})
}

// ResetDynaFields will reset dynamic fields.
func (i *Logger) ResetDynaFields() {
	i.dynamicFields = make([]interface{}, 2)
}
