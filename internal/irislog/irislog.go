// Package irislog defines custom context logger wrapped around rs/zerolog
package irislog

import (
	"context"
	"fmt"
	"os"

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
	// Can also be replaced by a custom configured one using Set(*Logger).
	StdLogger *IrisLogger
	cfg       Config
	// NOTE: future reference we also want to configure a log Config accessed for users using configparser.
)

// IrisLogger defines a default logger for iris that wraps around rs/zerolog.
type IrisLogger struct {
	Level      int
	Version    string
	Revision   string
	StdLog     zerolog.Logger
	ErrLog     zerolog.Logger
	dynafields []interface{}
}

// Config defines config for IrisLogger.
type Config struct {
	name       string
	level      int
	stfields   []interface{}
	configured bool
}

type ctxKey string

// String is human readable representation of a context key.
func (c ctxKey) String() string {
	return "mw-" + string(c)
}

func cfgSetup(name string, stfields []interface{}) {
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

// AddDynaFields acts as a receiver instance that add given key-value pairs to current logger.
func (i *IrisLogger) AddDynaFields(key, value interface{}) {
	i.dynafields = append(i.dynafields, []interface{}{key, value})
}

// ResetDynaFields will reset dynamic fields.
func (i *IrisLogger) ResetDynaFields() {
	i.dynafields = make([]interface{}, 2)
}

func init() {
	StdLogger = NewLogger(Debug, "irislog")
}

// CtxLogger returns a logger stored in the context provided by the arguments.
// We want to avoid runtime error, thus get the logger from context by using logger context keys.
func CtxLogger(ctx context.Context) (logger *IrisLogger, ok bool) {
	l, ok := ctx.Value(irisCtxKeys).(*IrisLogger)
	return l, ok
}

// NewLogger creates a new Logger.
// if static fields are provided those values will be defined.
// the default static fields for each new built instance if they aren't configured yet.
func NewLogger(level int, name string, stfields ...interface{}) *IrisLogger {
	if level < Disabled || level > Error {
		level = Info
	}

	stdl := zerolog.New(os.Stdout).With().Timestamp().Logger()
	errl := zerolog.New(os.Stderr).With().Timestamp().Logger()

	setupZerologLevel(&stdl, level)
	setupZerologLevel(&errl, level)

	i := &IrisLogger{
		Level:  level,
		StdLog: stdl,
		ErrLog: errl,
	}

	// NOTE: Possible workaround is to create a separate Config struct for log?
	if len(stfields) > 1 && !cfg.configured {
		cfgSetup(name, stfields)
		StdLogger = i
	}

	return i
}

// Set StdLogger to user's defined IrisLogger.
func Set(i *IrisLogger) {
	StdLogger = i
}

// Set chains NewLogger with StdLogger to create a new logger.
// l := irislog.NewLogger(irislog.Debug, "name", "version", "revision").Set().
func (i *IrisLogger) Set() *IrisLogger {
	StdLogger = i
	return StdLogger
}

// InCtx returns a copy of context that also includes a configured logger.
func InCtx(ctx context.Context, fields ...string) context.Context {
	l, _ := FromCtx(ctx)
	if len(fields) > 0 {
		l.SetDynaFields(fields)
	}
	return context.WithValue(ctx, irisCtxKeys, l)
}

// FromCtx returns current logger in context.
// If there isn't one then returns a new one with given cfg values.
func FromCtx(ctx context.Context) (i *IrisLogger, fresh bool) {
	l, ok := ctx.Value(irisCtxKeys).(IrisLogger)
	if !ok {
		return NewLogger(cfg.level, cfg.name), true
	}

	return &l, false
}

// Debug logs.
func (i IrisLogger) Debug(meta ...interface{}) {
	if len(meta) > 0 {
		i.debugf(stringify(meta[0]), meta[1:])
	}
}

func (i IrisLogger) debugf(message string, fields []interface{}) {
	if i.Level > Debug {
		return
	}
	ie := i.StdLog.Info()
	appendKeyValues(ie, i.dynafields, fields)
	ie.Msg(message)
}

// Info logs.
func (i IrisLogger) Info(meta ...interface{}) {
	if len(meta) > 0 {
		i.infof(stringify(meta[0]), meta[1:])
	}
}

func (i IrisLogger) infof(message string, fields []interface{}) {
	if i.Level > Info {
		return
	}

	ie := i.StdLog.Info()
	appendKeyValues(ie, i.dynafields, fields)
	ie.Msg(message)
}

// Warn logs.
func (i IrisLogger) Warn(meta ...interface{}) {
	if len(meta) > 0 {
		i.warnf(stringify(meta[0]), meta[1:])
	}
}

func (i IrisLogger) warnf(message string, fields []interface{}) {
	if i.Level > Warn {
		return
	}

	ie := i.StdLog.Warn()
	appendKeyValues(ie, i.dynafields, fields)
	ie.Msg(message)
}

// Error logs.
func (i IrisLogger) Error(err error, meta ...interface{}) {
	if len(meta) > 0 {
		i.errorf(err, stringify(meta[0]), meta[1:])
		return
	}
	i.errorf(err, "", nil)
}

func (i IrisLogger) errorf(err error, message string, fields []interface{}) {
	ie := i.ErrLog.Error()
	appendKeyValues(ie, i.dynafields, fields)
	ie.Err(err)
	ie.Msg(message)
}

func appendKeyValues(le *zerolog.Event, dynafields, fields []interface{}) {
	if cfg.name != "" {
		le.Str("name", cfg.name)
	}

	fs := make(map[string]interface{})

	// TODO: static key-value should be cached?
	if len(fields) > 1 {
		for i := 0; i < len(fields)-1; i++ {
			if fields[i] != nil && fields[i+1] != nil {
				k := stringify(fields[i])
				fs[k] = fields[i+1]
				i++
			}
		}

		if len(dynafields) > 1 {
			for i := 0; i < len(dynafields)-1; i++ {
				if dynafields[i] != nil && dynafields[i+1] != nil {
					k := stringify(dynafields[i])
					fs[k] = dynafields[i+1]
					i++
				}
			}
		}

		if len(cfg.stfields) > 1 {
			for i := 0; i < len(cfg.stfields)-1; i++ {
				if cfg.stfields[i] != nil && cfg.stfields[i+1] != nil {
					k := stringify(cfg.stfields[i])
					fs[k] = cfg.stfields[i+1]
					i++
				}
			}
		}
	}

	le.Fields(fs)
}

func stringify(val interface{}) string {
	switch v := val.(type) {
	case nil:
		return fmt.Sprintf("%v", v)
	case int:
		return fmt.Sprintf("%d", v)
	case bool:
		return fmt.Sprintf("%t", v)
	case string:
		return v
	default:
		return fmt.Sprintf("%+v", v)
	}
}

// UpdateLogLevel updates log level.
func (i *IrisLogger) UpdateLogLevel(level int) {
	// don't downgrade if current is Error
	current := Error
	i.Info("Log level updated", "", "log level", level)
	i.Level = current
	if level < Disabled || level > Error {
		i.Level = level
		setupZerologLevel(&i.StdLog, level)
		setupZerologLevel(&i.ErrLog, level)
	}
}

// this will setup correct log level for our zerolog.
func setupZerologLevel(l *zerolog.Logger, level int) {
	switch level {
	case -1:
		l.Level(zerolog.Disabled)
	case 0:
		l.Level(zerolog.DebugLevel)
	case 1:
		l.Level(zerolog.InfoLevel)
	case 2:
		l.Level(zerolog.WarnLevel)
	case 3:
		l.Level(zerolog.ErrorLevel)
	default:
		l.Level(zerolog.DebugLevel)
	}
}
