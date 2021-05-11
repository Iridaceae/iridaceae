package log

import (
	"bytes"
	"errors"
	"testing"

	"github.com/rs/zerolog"

	"github.com/stretchr/testify/assert"
)

func TestScLevelEncoder(t *testing.T) {
	original := zerolog.LevelFieldMarshalFunc
	zerolog.LevelFieldMarshalFunc = ScLevelEncoder()
	defer func() {
		zerolog.LevelFieldMarshalFunc = original
	}()

	tests := []struct {
		name    string
		msg     string
		want    string
		logFunc func()
	}{
		{"debug logs", "test", `{"level":"DEBUG","message":"test"}` + "\n", func() { Debug().Msg("test") }},
		{"info logs", "test", `{"level":"INFO","message":"test"}` + "\n", func() { Info().Msg("test") }},
		{"warn logs", "test", `{"level":"WARN","message":"test"}` + "\n", func() { Warn().Msg("test") }},
		{"error logs", "test", `{"level":"ERROR","error":"test"}` + "\n", func() { Error(errors.New("test")).Msg("") }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			L = NewZ(zerolog.New(out))
			tt.logFunc()
			if got, want := decodeIfBinaryToString(out.Bytes()), tt.want; got != want {
				t.Errorf("invalid log output:\ngot: %v\nwant: %v", got, want)
			}
			out.Reset()
		})
	}
}

func TestScCallerEncoder(t *testing.T) {

}

func TestTrimmedPath(t *testing.T) {
	tests := []struct {
		name string
		want string
		got  string
	}{
		{"valid path", "test/test.go", TrimmedPath("test/test.go")},
		{"longer path", "test/test.go", TrimmedPath("hello/world/test/test.go")},
		{"invalid path", "test-test.go", TrimmedPath("test-test.go")}, // invalid path shouldn't trim anything
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.got)
		})
	}
}
