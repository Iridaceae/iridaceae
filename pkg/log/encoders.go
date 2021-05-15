package log

import (
	"strconv"
	"strings"

	"github.com/rs/zerolog"
)

func LevelEncoder() func(l zerolog.Level) string {
	return func(l zerolog.Level) string {
		return strings.ToUpper(l.String())
	}
}

func CallerEncoder() func(file string, line int) string {
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
