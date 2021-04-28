package irislog

import (
	"github.com/rs/zerolog"
)

const (
	irisCtxKeys ctxKey = "irislog"
	Disabled           = -1
	Debug              = iota
	Info
	Warn
	Error
)

var (
	// StdLogger Logger can be used right out of the box.
	// Can also be replaced by a custom configured one using Set(*Logger)
	StdLogger *IrisLogger
	cfg       config
	// NOTE: future reference we also want to configure a log config accessed for users using configparser
)

type IrisLogger struct {
	Level      int
	Version    string
	Revision   string
	StdLog     zerolog.Logger
	ErrLog     zerolog.Logger
	dynafields []interface{}
}

type config struct {
	name       string
	level      int
	stfields   []interface{}
	configured bool
}

type ctxKey string

// String is human readable representation of a context key
func (c ctxKey) String() string {
	return "mw-" + string(c)
}

func setup(name string, stfields []interface{}) {
	if cfg.configured {
		return
	}

	cfg.name = name
	cfg.stfields = append(cfg.stfields, stfields...)
	cfg.configured = true
}

// SetDynaFields acts as a receiver instance that will always append these key-value pairs to the output.
func (i *IrisLogger) SetDynaFields(dynafields ...interface{}) {
	i.dynafields = make([]interface{}, 2)
	i.dynafields = append(i.dynafields, dynafields...)
}

// AddDynaFields acts as a receiver instance that add given key-value pairs to current logger
func (i *IrisLogger) AddDynaFields(key, value interface{}) {
	i.dynafields = append(i.dynafields, []interface{}{key, value})
}

// ResetDynaFields will reset dynamic fields
func (i *IrisLogger) ResetDynaFields() {
	i.dynafields = make([]interface{}, 2)
}
