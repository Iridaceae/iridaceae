// Package stlog provides a structured context logger.
package stlog

import (
	"fmt"
	"os"

	"github.com/Iridaceae/iridaceae/pkg/util"

	"github.com/rs/zerolog"
)

const (
	logCtxKey LogCtxKey = "certlog"
	Disabled            = -1
	Debug               = iota
	Info
	Warn
	Error
)

// Defaults can be used right out of the box.
// Can also be replaced by a custom configured one using Set(*Logger).
var Defaults *Logger

func init() {
	Defaults = NewLogger(Debug, logCtxKey.String())
	Defaults.SetDynaFields("revision", util.GetRevision())
}

// NewLogger creates a new Logger.
// if static fields are provided those values will be defined.
// the default static fields for each new built instance if they aren't configured yet.
func NewLogger(level int, name string, staticFields ...interface{}) *Logger {
	if level < Disabled || level > Error {
		level = Info
	}

	stdLog := zerolog.New(os.Stdout).
		With().
		Timestamp().
		Logger()
	errLog := zerolog.New(os.Stderr).
		With().
		Timestamp().
		Logger()

	setLogLevel(&stdLog, level)
	setLogLevel(&errLog, level)

	i := &Logger{
		Level:  level,
		StdLog: stdLog,
		ErrLog: errLog,
	}

	// NOTE: Possible workaround is to create a separate LogConfig struct for log?
	if len(staticFields) > 1 && !cfg.configured {
		cfgSetup(name, staticFields)
		Defaults = i
	}

	return i
}

// Set Defaults to user's defined Logger.
func Set(i *Logger) {
	Defaults = i
}

// Set chains NewLogger with Defaults to create a new logger.
// 		l := log.NewLogger(log.Debug, "name", "version", revision).Set().
func (i *Logger) Set() *Logger {
	Defaults = i
	return Defaults
}

// Debug logs.
func (i Logger) Debug(meta ...interface{}) {
	if len(meta) > 0 {
		i.debugf(stringify(meta[0]), meta[1:])
	}
}

func (i Logger) debugf(message string, fields []interface{}) {
	if i.Level > Debug {
		return
	}
	ie := i.StdLog.Info()
	appendKeyValues(ie, i.dynamicFields, fields)
	ie.Msg(message)
}

// Info logs.
func (i Logger) Info(meta ...interface{}) {
	if len(meta) > 0 {
		i.infof(stringify(meta[0]), meta[1:])
	}
}

func (i Logger) infof(message string, fields []interface{}) {
	if i.Level > Info {
		return
	}

	ie := i.StdLog.Info()
	appendKeyValues(ie, i.dynamicFields, fields)
	ie.Msg(message)
}

// Warn logs.
func (i Logger) Warn(meta ...interface{}) {
	if len(meta) > 0 {
		i.warnf(stringify(meta[0]), meta[1:])
	}
}

func (i Logger) warnf(message string, fields []interface{}) {
	if i.Level > Warn {
		return
	}

	ie := i.StdLog.Warn()
	appendKeyValues(ie, i.dynamicFields, fields)
	ie.Msg(message)
}

// Error logs.
func (i Logger) Error(err error, meta ...interface{}) {
	if len(meta) > 0 {
		i.errorf(err, stringify(meta[0]), meta[1:])
		return
	}
	i.errorf(err, "", nil)
}

func (i Logger) errorf(err error, message string, fields []interface{}) {
	ie := i.ErrLog.Error()
	appendKeyValues(ie, i.dynamicFields, fields)
	ie.Err(err)
	ie.Msg(message)
}

func appendKeyValues(le *zerolog.Event, dynamicFields, staticFields []interface{}) {
	if cfg.name != "" {
		le.Str("name", cfg.name)
	}

	fs := make(map[string]interface{})

	// TODO: static key-value should be cached?
	if len(staticFields) > 1 {
		for i := 0; i < len(staticFields)-1; i++ {
			if staticFields[i] != nil && staticFields[i+1] != nil {
				k := stringify(staticFields[i])
				fs[k] = staticFields[i+1]
				i++
			}
		}

		if len(dynamicFields) > 1 {
			for i := 0; i < len(dynamicFields)-1; i++ {
				if dynamicFields[i] != nil && dynamicFields[i+1] != nil {
					k := stringify(dynamicFields[i])
					fs[k] = dynamicFields[i+1]
					i++
				}
			}
		}

		if len(cfg.staticFields) > 1 {
			for i := 0; i < len(cfg.staticFields)-1; i++ {
				if cfg.staticFields[i] != nil && cfg.staticFields[i+1] != nil {
					k := stringify(cfg.staticFields[i])
					fs[k] = cfg.staticFields[i+1]
					i++
				}
			}
		}
	}
	le.Fields(fs)
}

// UpdateLogLevel updates log level.
func (i *Logger) UpdateLogLevel(level int) {
	// don't downgrade if current is Error
	current := Error
	i.Info("Log level updated", "", "log level", level)
	i.Level = current
	if level < Disabled || level > Error {
		i.Level = level
		setLogLevel(&i.StdLog, level)
		setLogLevel(&i.ErrLog, level)
	}
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

func setLogLevel(l *zerolog.Logger, level int) {
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
