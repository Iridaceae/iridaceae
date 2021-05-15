package rosetta

import (
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasPrefix(t *testing.T) {
	testPrefixFunc := func(msg string, prefix string, ignoreCase bool, ok bool) {
		_, k := hasPrefix(msg, prefix, ignoreCase)
		assert.Equal(t, k, ok)
	}

	tests := []struct {
		name       string
		expected   bool
		msg        string
		prefix     string
		prefixFunc func(s string, prefix string, ignoreCase bool, ok bool)
	}{
		{"doesn't have prefix contain in given string", false, "hello world", "!", testPrefixFunc},
		{"has prefix", true, "!hello world", "!", testPrefixFunc},
		{"has complex prefix", true, "!ir hello world", "!ir", testPrefixFunc},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prefixFunc(tt.msg, tt.prefix, true, tt.expected)
		})
	}
}

func TestTrimPreSuffix(t *testing.T) {
	s := "'hello world'"
	preSuffix := "'"
	o := trimPreSuffix(s, preSuffix)
	assert.Equal(t, "hello world", o)
}

func TestArrayContains(t *testing.T) {
	t.Run("ignore case", func(t *testing.T) {
		arr := []string{"1", "2", "3"}
		contained := "test"
		ok := arrayContains(arr, contained, false)
		assert.False(t, ok)
	})
	t.Run("case sensitive", func(t *testing.T) {
		assert.True(t, arrayContains([]string{"TEST", "2", "3"}, "test", true))
	})
}

func TestClearMap(t *testing.T) {
	m := sync.Map{}
	m.Store("test", 1)
	clearMap(&m)
	_, ok := m.Load("test")
	assert.False(t, ok)
}

func TestGetErrorType(t *testing.T) {
	tests := []struct {
		name   string
		e      error
		output string
	}{
		{"invalid error return empty", errors.New(""), getErrorTypeName(-1)},
		{"valid error", ErrGuildPrefixGetter, getErrorTypeName(0)},
		{"error while get channel", ErrGetChannel, getErrorTypeName(1)},
		{"error while get guild", ErrGetGuild, getErrorTypeName(2)},
		{"error cannot find command", ErrCommandNotFound, getErrorTypeName(3)},
		{"error command not executable in dm", ErrNotExecutableInDMs, getErrorTypeName(4)},
		{"error middleware", ErrMiddleware, getErrorTypeName(5)},
		{"error command exec", ErrCommandExec, getErrorTypeName(6)},
		{"error delete command message", ErrDeleteCommandMessage, getErrorTypeName(7)},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%s-%d", tt.name, i), func(t *testing.T) {
			assert.Equal(t, tt.e.Error(), tt.output)
		})
	}
}
