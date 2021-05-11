package log

import (
	"testing"

	"github.com/Iridaceae/iridaceae/internal/testutils/cbor"

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
