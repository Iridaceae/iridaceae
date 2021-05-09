package sclog

import (
	"bytes"
	"testing"

	"github.com/rs/zerolog"

	"github.com/Iridaceae/iridaceae/internal/cbor"

	"github.com/stretchr/testify/assert"
)

func decodeIfBinaryToString(in []byte) string {
	return cbor.DecodeIfBinaryToString(in)
}

func TestAddGlobalFields(t *testing.T) {
	AddGlobalFields("test")
	assert.Equal(t, 1, len(GetGlobalFields()))
}

func TestSetGlobalFields(t *testing.T) {
	SetGlobalFields([]string{"test", "test2"})
	assert.Equal(t, 2, len(GetGlobalFields()))
}

func TestClearGlobalFields(t *testing.T) {
	ClearGlobalFields()
	assert.Equal(t, 0, len(GetGlobalFields()))
}

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
		logFunc func(s *ScLogger)
	}{
		{"debug logs", "test", `{"level":"DEBUG","message":"test"}` + "\n", func(s *ScLogger) { s.Debug("test") }},
		{"info logs", "test", `{"level":"INFO","message":"test"}` + "\n", func(s *ScLogger) { s.Info("test") }},
		{"warn logs", "test", `{"level":"WARN","message":"test"}` + "\n", func(s *ScLogger) { s.Warn("test") }},
		{"error logs", "test", `{"level":"ERROR","message":"test"}` + "\n", func(s *ScLogger) { s.Error("test") }},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			log := SetZ(zerolog.New(out))
			tt.logFunc(log)
			if got, want := decodeIfBinaryToString(out.Bytes()), tt.want; got != want {
				t.Errorf("invalid log output:\ngot: %v\nwant: %v", got, want)
			}
			out.Reset()
		})
	}
}

func TestTrimmedPath(t *testing.T) {

}
