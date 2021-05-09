// Package log provides a front facing access to sclog.
package log

import (
	"github.com/Iridaceae/iridaceae/pkg/sclog"
	"github.com/Iridaceae/iridaceae/pkg/util"

	"github.com/rs/zerolog"
)

// Logger is a global Logger that implements rs/zerlog.
// this can also be replaced by a custom configured one using SetZ(*zerolog.Logger).
var Logger = New()

// New duplicates scLog.New() with our custom fields.
func New() *sclog.ScLogger {
	sclog.New()
	sclog.Mapper().Set("revision", util.GetRevision())
	sclog.Mapper().Set("version", util.GetVersion())
	sclog.SetGlobalFields([]string{"revision", "version"})
	return sclog.Log()
}

// SetZ duplicates global logger and returns a new custom user-defined zerolog.
func SetZ(l zerolog.Logger) *sclog.ScLogger {
	Logger = sclog.SetZ(l)
	return Logger
}

// Trace logs.
func Trace(args ...interface{}) {
	Logger.Trace(args...)
}

// Tracef format logs.
func Tracef(format string, args ...interface{}) {
	Logger.Tracef(format, args...)
}

// Debug logs.
func Debug(args ...interface{}) {
	Logger.Debug(args...)
}

// Debugf format logs.
func Debugf(format string, args ...interface{}) {
	Logger.Debugf(format, args...)
}

// Info logs.
func Info(args ...interface{}) {
	Logger.Info(args...)
}

// Infof format logs.
func Infof(format string, args ...interface{}) {
	Logger.Infof(format, args...)
}

// Warn logs.
func Warn(args ...interface{}) {
	Logger.Warn(args...)
}

// Warnf format logs.
func Warnf(format string, args ...interface{}) {
	Logger.Warnf(format, args...)
}

// Error logs.
func Error(args ...interface{}) {
	Logger.Error(args...)
}

// Errorf format logs.
func Errorf(format string, args ...interface{}) {
	Logger.Errorf(format, args...)
}

// Fatal logs.
func Fatal(args ...interface{}) {
	Logger.Fatal(args...)
}

// Fatalf format logs.
func Fatalf(format string, args ...interface{}) {
	Logger.Fatalf(format, args...)
}

// Panic logs.
func Panic(args ...interface{}) {
	Logger.Panic(args...)
}

// Panicf format logs.
func Panicf(format string, args ...interface{}) {
	Logger.Panicf(format, args...)
}

// Log logs.
func Log(args ...interface{}) {
	Logger.Log(args...)
}

// Logf format logs.
func Logf(format string, args ...interface{}) {
	Logger.Logf(format, args...)
}

// Print logs.
func Print(args ...interface{}) {
	Logger.Print(args...)
}

// Printf format logs.
func Printf(format string, args ...interface{}) {
	Logger.Printf(format, args...)
}
