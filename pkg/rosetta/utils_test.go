package rosetta

import (
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/bwmarrin/discordgo"

	"github.com/stretchr/testify/assert"
)

var TestSession = &discordgo.Session{State: &discordgo.State{Ready: discordgo.Ready{User: &discordgo.User{ID: "123465879"}}}}

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
		tarr := []string{"1", "2", "3"}
		contained := "test"
		ok := arrayContains(tarr, contained, false)
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
		{"invalid error return empty", errors.New(""), getErrorTypeName(0)},
		{"valid error", ErrGuildPrefixGetter, getErrorTypeName(1)},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%s-%d", tt.name, i), func(t *testing.T) {
			assert.Equal(t, tt.e.Error(), tt.output)
		})
	}
}
