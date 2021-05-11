package log

import (
	"bytes"
	"errors"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

// hooks fn.
var (
	simpleHook = zerolog.HookFunc(func(e *zerolog.Event, level zerolog.Level, message string) {
		e.Bool("hasLevel", level != zerolog.NoLevel)
	})
)

func TestNew(t *testing.T) {
	zlog := zerolog.New(os.Stdout).With().Timestamp().Logger()
	zhook := zlog.Hook(MapperHook{})
	actual := &zhook

	NewZ(zlog)
	expected := L.log
	assert.Equal(t, expected, actual)
	assert.Equal(t, expected, Z())
}

func TestEvents(t *testing.T) {
	tests := []struct {
		name string
		want string
		log  func(l zerolog.Logger)
	}{
		{"log level", `{"hasLevel":false,"message":"test"}` + "\n", func(l zerolog.Logger) {
			l = l.Hook(simpleHook)
			L.log = &l
			Log().Msg("test")
		}},
		{"print", `{"level":"DEBUG","hasLevel":true,"message":"test"}` + "\n", func(l zerolog.Logger) {
			l = l.Hook(simpleHook)
			L.log = &l
			Print("test")
		}},
		{"printf", `{"level":"DEBUG","hasLevel":true,"message":"test hello"}` + "\n", func(l zerolog.Logger) {
			l = l.Hook(simpleHook)
			L.log = &l
			Printf("test %s", "hello")
		}},
		{"trace level", `{"level":"TRACE","hasLevel":true,"message":"test"}` + "\n", func(l zerolog.Logger) {
			l = l.Hook(simpleHook)
			L.log = &l
			Trace().Msg("test")
		}},
		{"debug level", `{"level":"DEBUG","hasLevel":true,"message":"test"}` + "\n", func(l zerolog.Logger) {
			l = l.Hook(simpleHook)
			L.log = &l
			Debug().Msg("test")
		}},
		{"info level", `{"level":"INFO","hasLevel":true,"message":"test"}` + "\n", func(l zerolog.Logger) {
			l = l.Hook(simpleHook)
			L.log = &l
			Info().Msg("test")
		}},
		{"warn level", `{"level":"WARN","hasLevel":true,"message":"test"}` + "\n", func(l zerolog.Logger) {
			l = l.Hook(simpleHook)
			L.log = &l
			Warn().Msg("test")
		}},
		{"error level", `{"level":"ERROR","error":"test error","hasLevel":true,"message":"test"}` + "\n", func(l zerolog.Logger) {
			l = l.Hook(simpleHook)
			L.log = &l
			err := errors.New("test error")
			Error(err).Msg("test")
		}},
		// TODO : fatal and panic tests
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			L = NewZ(zerolog.New(out))
			tt.log(*L.log)
			if got, want := decodeIfBinaryToString(out.Bytes()), tt.want; got != want {
				t.Errorf("invalid log output:\ngot: %v\nwant: %v", got, want)
			}
		})
	}
}
