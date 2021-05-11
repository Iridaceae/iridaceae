// Package log provides a structured context logger.
package log

import (
	"os"
	"time"

	"github.com/rs/zerolog"

	"github.com/rs/zerolog/log"
)

var L *Logger

func init() {
	L = New()
}

// Logger defines a default logger for that wraps around rs/zerolog.
type Logger struct {
	log *zerolog.Logger
}

// NewZ creates a Logger from user-defined zerolog.
func NewZ(l zerolog.Logger) *Logger {
	log.Logger = l.Hook(MapperHook{})
	ResetGlobalStorage()
	ClearGlobalFields()

	L = &Logger{log: &log.Logger}
	return L
}

// New returns a default Logger wrapped around zerolog.
func New() *Logger {
	l := zerolog.New(os.Stdout).With().Timestamp().Caller().Logger().Hook(MapperHook{})
	ResetGlobalStorage()
	ClearGlobalFields()

	zerolog.CallerFieldName = "source"
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.LevelFieldMarshalFunc = ScLevelEncoder()
	zerolog.CallerMarshalFunc = ScCallerEncoder()

	L = &Logger{log: &l}
	return L
}

// Z returns internal zerolog.Logger of our global logger.
func Z() *zerolog.Logger {
	return L.log
}

func Trace() *zerolog.Event {
	return L.log.Trace()
}

func Debug() *zerolog.Event {
	return L.log.Debug()
}

func Info() *zerolog.Event {
	return L.log.Info()
}

func Warn() *zerolog.Event {
	return L.log.Warn()
}

func Error(err error) *zerolog.Event {
	return L.log.Error().Err(err)
}

func Fatal(err error) *zerolog.Event {
	return L.log.Fatal().Err(err)
}

func Panic() *zerolog.Event {
	return L.log.Panic()
}

func Log() *zerolog.Event {
	return L.log.Log()
}

func Print(args ...interface{}) {
	L.log.Print(args...)
}

func Printf(format string, args ...interface{}) {
	L.log.Printf(format, args...)
}
