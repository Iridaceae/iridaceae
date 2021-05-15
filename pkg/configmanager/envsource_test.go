package configmanager

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var envSource = &EnvSource{}

func TestEnvSource_GetValue(t *testing.T) {
	t.Run("get existed envars", func(t *testing.T) {
		setupConfigTest(t)
		var key, value string
		key = "TEST_OPTION1_OPTION2"
		value = "test_option1_option2"
		opt := "test.option1.option2"
		createTestEnvVars(t, key, value)
		v, err := envSource.GetValue(opt)
		assert.Equal(t, v, value)
		assert.Nil(t, err)
	})

	t.Run("get nil envars", func(t *testing.T) {
		setupConfigTest(t)
		var key, value string
		key = "TEST_NIL"
		value = ""
		opt := "test.nil"
		createTestEnvVars(t, key, value)
		v, err := envSource.GetValue(opt)
		assert.Nil(t, v)
		assert.ErrorIs(t, err, ErrEmptyValue)
	})

	t.Run("test regex with uncorrect envars", func(t *testing.T) {
		setupConfigTest(t)
		var key, value string
		key = "TEST_FORMAT1"
		value = "test_format1"
		opt := []string{"test.format1+format2", "test-format1.format2", "test&format1"}
		createTestEnvVars(t, key, value)
		for _, val := range opt {
			v, err := envSource.GetValue(val)
			assert.NotEqual(t, v, value)
			assert.ErrorIs(t, err, ErrInvalidFormat)
		}
	})
}

func TestEnvSource_Name(t *testing.T) {
	t.Run("get name", func(t *testing.T) {
		setupConfigTest(t)
		assert.Equal(t, envSource.Name(), "ENV")
	})
}
