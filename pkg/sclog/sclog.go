// Package sclog provides a structured context logger.
package sclog

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

// SetZ creates a ScLogger from user-defined zerolog.
func SetZ(l zerolog.Logger) *ScLogger {
	l = l.Hook(MapperHook{})
	ResetGlobalStorage()
	ClearGlobalFields()

	scLog = &ScLogger{log: &l}
	return scLog
}

// New returns a default ScLogger wrapped around zerolog.
func New() *ScLogger {
	l := zerolog.New(os.Stdout).With().Timestamp().Caller().Logger().Hook(MapperHook{})
	ResetGlobalStorage()
	ClearGlobalFields()

	zerolog.CallerFieldName = "source"
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.LevelFieldMarshalFunc = ScLevelEncoder()
	zerolog.CallerMarshalFunc = ScCallerEncoder()

	scLog = &ScLogger{log: &l}
	return scLog
}

// Log returns our internal logger.
func Log() *ScLogger {
	return scLog
}

// Z returns internal zerolog.Logger.
func Z() *zerolog.Logger {
	return scLog.log
}
func (s *ScLogger) Trace(args ...interface{}) {
	s.log.Trace().Msg(message("", args...))
}

func (s *ScLogger) Tracef(format string, args ...interface{}) {
	s.log.Trace().Msg(message(format, args...))
}

func (s *ScLogger) Debug(args ...interface{}) {
	s.log.Debug().Msg(message("", args...))
}

func (s *ScLogger) Debugf(format string, args ...interface{}) {
	s.log.Debug().Msg(message(format, args...))
}

func (s *ScLogger) Info(args ...interface{}) {
	s.log.Info().Msg(message("", args...))
}

func (s *ScLogger) Infof(format string, args ...interface{}) {
	s.log.Info().Msg(message(format, args...))
}

func (s *ScLogger) Warn(args ...interface{}) {
	s.log.Warn().Msg(message("", args...))
}

func (s *ScLogger) Warnf(format string, args ...interface{}) {
	s.log.Warn().Msg(message(format, args...))
}

func (s *ScLogger) Error(args ...interface{}) {
	s.log.Error().Msg(message("", args...))
}

func (s *ScLogger) Errorf(format string, args ...interface{}) {
	s.log.Error().Msg(message(format, args...))
}

func (s *ScLogger) Fatal(args ...interface{}) {
	s.log.Fatal().Msg(message("", args...))
}

func (s *ScLogger) Fatalf(format string, args ...interface{}) {
	s.log.Fatal().Msg(message(format, args...))
}

func (s *ScLogger) Panic(args ...interface{}) {
	s.log.Panic().Msg(message("", args...))
}

func (s *ScLogger) Panicf(format string, args ...interface{}) {
	s.log.Panic().Msg(message(format, args...))
}

func (s *ScLogger) Log(args ...interface{}) {
	s.log.Log().Msg(message("", args...))
}

func (s *ScLogger) Logf(format string, args ...interface{}) {
	s.log.Log().Msg(message(format, args...))
}

func (s *ScLogger) Print(args ...interface{}) {
	s.log.Print(args...)
}

func (s *ScLogger) Printf(format string, args ...interface{}) {
	s.log.Printf(format, args...)
}
