package sclog

import (
	"strconv"
	"strings"

	"github.com/rs/zerolog"
)

var _fields = make([]string, 0)

func SetGlobalFields(fields []string) {
	_fields = fields
}

func AddGlobalFields(field string) {
	_fields = append(_fields, field)
}

func GetGlobalFields() []string {
	return _fields
}

func ClearGlobalFields() {
	_fields = make([]string, 0)
}

func ScLevelEncoder() func(l zerolog.Level) string {
	return func(l zerolog.Level) string {
		return strings.ToUpper(l.String())
	}
}

func ScCallerEncoder() func(file string, line int) string {
	return func(file string, line int) string {
		return TrimmedPath(file) + ":" + strconv.Itoa(line)
	}
}

func TrimmedPath(file string) string {
	idx := strings.LastIndexByte(file, '/')
	if idx == -1 {
		return file
	}
	idx = strings.LastIndexByte(file[:idx], '/')
	if idx == -1 {
		return file
	}
	return file[idx+1:]
}
