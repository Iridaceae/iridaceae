package sclog

import (
	"fmt"

	"github.com/rs/zerolog"
)

var scLog *ScLogger

// ScLogger defines a default logger for that wraps around rs/zerolog.
type ScLogger struct {
	log *zerolog.Logger
}

func message(format string, args ...interface{}) string {
	msg := format
	if msg == "" && len(args) > 0 {
		msg = fmt.Sprint(args...)
	} else if msg != "" && len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	}
	return msg
}
