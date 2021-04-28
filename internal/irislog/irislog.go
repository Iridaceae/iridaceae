// Package irislog defines custom log wrapper around rs/zerolog
package irislog

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

var (
	// Defaults logger can be used right out of the box.
	// Can also be replaced by a custom configured one using Set(*Logger)
	Defaults *IrisLogger
)

func init() {
	Defaults = NewLogger(Debug, "irislog")
}

// NewLogger creates a new logger.
// if static fields are provided those values will be defined.
// the default static fields for each new built instance if they aren't configured yet.
func NewLogger(level int, name string, stfields ...interface{}) *IrisLogger {
	if level < Disabled || level > Error {
		level = Info
	}

	stdl := zerolog.New(os.Stdout).With().Timestamp().Logger()
	errl := zerolog.New(os.Stderr).With().Timestamp().Logger()

	setLogLevel(&stdl, level)
	setLogLevel(&errl, level)

	i := &IrisLogger{
		Level:  level,
		Name:   name,
		StdLog: stdl,
		ErrLog: errl,
	}

	// NOTE: Possible workaround is to create a separate config struct for log?
	if len(stfields) > 1 && !logCfg.manager.Options["irislog.configured"].GetBool() {
		setup(level, stfields)
		Defaults = i
	}

	return i
}

// NewDevLogger creates development logger.
// Pretty logging for development mode.
func NewDevLogger(level int, name string, stfields ...interface{}) *IrisLogger {
	if level < Disabled || level > Error {
		level = Info
	}

	stdl := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout})
	errl := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr})

	setLogLevel(&stdl, level)
	setLogLevel(&errl, level)

	i := &IrisLogger{
		Level:  level,
		Name:   name,
		StdLog: stdl,
		ErrLog: errl,
	}

	if len(stfields) > 1 && !logCfg.manager.Options["irislog.configured"].GetBool() {
		setup(level, stfields)
		Defaults = i
	}

	return i
}

// Set setup base logger
func Set(i *IrisLogger) {
	Defaults = i
	logCfg.logger = i
}

// Set setup base logger
// Can be used to chain with NewLogger to create a new logger
// ```logger := irislog.NewLogger(irislog.Debug, "name", "version", "revision").Set()```.
func (i *IrisLogger) Set() *IrisLogger {
	Defaults = i
	logCfg.logger = Defaults
	return Defaults
}

func (i IrisLogger) Debug(meta ...interface{}) {
	if len(meta) > 0 {
		i.debugf(stringify(meta[0]), meta[1:len(meta)])
	}
}

func (i IrisLogger) debugf(message string, fields []interface{}) {
	if i.Level > Debug {
		return
	}
	ie := i.StdLog.Info()
	appendKeyValues(ie, fields)
	ie.Msg(message)
}

func (i IrisLogger) Info(meta ...interface{}) {
	if len(meta) > 0 {
		i.infof(stringify(meta[0]), meta[1:len(meta)])
	}
}

func (i IrisLogger) infof(message string, fields []interface{}) {
	if i.Level > Info {
		return
	}

	ie := i.StdLog.Info()
	appendKeyValues(ie, fields)
	ie.Msg(message)
}

func (i IrisLogger) Warn(meta ...interface{}) {
	if len(meta) > 0 {
		i.warnf(stringify(meta[0]), meta[1:len(meta)])
	}
}

func (i IrisLogger) warnf(message string, fields []interface{}) {
	if i.Level > Warn {
		return
	}

	ie := i.StdLog.Warn()
	appendKeyValues(ie, fields)
	ie.Msg(message)
}

func (i IrisLogger) Error(err error, meta ...interface{}) {
	if len(meta) > 0 {
		i.errorf(err, stringify(meta[0]), meta[1:len(meta)])
		return
	}
	i.errorf(err, "", nil)
}

func (i IrisLogger) errorf(err error, message string, fields []interface{}) {
	ie := i.ErrLog.Error()
	appendKeyValues(ie, fields)
	ie.Err(err)
	ie.Msg(message)
}

func appendKeyValues(le *zerolog.Event, fields []interface{}) {
	if logCfg.logger.Name != "" {
		le.Str("name", logCfg.logger.Name)
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

		if len(logCfg.stdFields) > 1 {
			for i := 0; i < len(logCfg.stdFields)-1; i++ {
				if logCfg.stdFields[i] != nil && logCfg.stdFields[i+1] != nil {
					k := stringify(logCfg.stdFields[i])
					fs[k] = logCfg.stdFields[i+1]
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

// UpdateLogLevel updates log level
func (i *IrisLogger) UpdateLogLevel(level int) {
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
