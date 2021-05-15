package configparser

import (
	"testing"
)

func TestEnvSource_GetValue(t *testing.T) {
	t.Run("get existed envars", func(t *testing.T) {
		setupConfigTest(t)
		var key, value string
		key = "TEST_OPTION1_OPTION2"
		value = "test_option1_option2"
		opt := "test.option1.option2"
		createTestEnvVars(t, key, value)
		v, err := testSource.GetValue(opt)
		cfgAssert.Equal(v, value)
		cfgAssert.Nil(err)
	})

	t.Run("get nil envars", func(t *testing.T) {
		setupConfigTest(t)
		var key, value string
		key = "TEST_NIL"
		value = ""
		opt := "test.nil"
		createTestEnvVars(t, key, value)
		v, err := testSource.GetValue(opt)
		cfgAssert.Nil(v)
		cfgAssert.ErrorIs(err, ErrEmptyValue)
	})

	t.Run("test regex with uncorrect envars", func(t *testing.T) {
		setupConfigTest(t)
		var key, value string
		key = "TEST_FORMAT1"
		value = "test_format1"
		opt := []string{"test.format1+format2", "test-format1.format2", "test&format1"}
		createTestEnvVars(t, key, value)
		for _, val := range opt {
			v, err := testSource.GetValue(val)
			cfgAssert.NotEqual(v, value)
			cfgAssert.ErrorIs(err, ErrInvalidFormat)
		}
	})
}

func TestEnvSource_Name(t *testing.T) {
	t.Run("get name", func(t *testing.T) {
		setupConfigTest(t)
		cfgAssert.Equal(testSource.Name(), "ENV")
	})
}
