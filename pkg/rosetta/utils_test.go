package rosetta

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasPrefix(t *testing.T) {
	t.Run("doesn't have prefix contain in given string", func(t *testing.T) {
		s := "hello world"
		prefs := []string{"!", "-"}
		ok, _ := hasPrefix(s, prefs, true)
		assert.False(t, ok)
	})
	t.Run("does have prefix in given string", func(t *testing.T) {
		s := "!hello world"
		prefs := []string{"!", "-"}
		ok, so := hasPrefix(s, prefs, true)
		assert.True(t, ok)
		assert.Equal(t, "hello world", so)
	})
}

func TestTrimPreSuffix(t *testing.T) {
	s := "'hello world'"
	preSuffix := "'"
	o := trimPreSuffix(s, preSuffix)
	assert.Equal(t, "hello world", o)
}

func TestArrayContains(t *testing.T) {
	tarr := []string{"1", "2", "3"}
	contained := "test"
	ok := arrayContains(tarr, contained, false)
	assert.False(t, ok)
}

func getEnvOrDefault(env, def string) string {
	v := os.Getenv(env)
	if v == "" && def != "" {
		v = def
	}
	return v
}
