package configparser

import (
	"testing"
)

func TestEnvSource_GetValue(t *testing.T) {
	t.Run("get existed envars", func(t *testing.T) {
		setup(t)
		var key, value string
		key = "TEST_OPTION1_OPTION2"
		value = "test_option1_option2"
		opt := "test.option1.option2"
		createEnvVar(t, key, value)
		v, err := testSource.GetValue(opt)
		defaultAssert.Equal(v, value)
		defaultAssert.Nil(err)
	})

	t.Run("get nil envars", func(t *testing.T) {
		setup(t)
		var key, value string
		key = "TEST_NIL"
		value = ""
		opt := "test.nil"
		createEnvVar(t, key, value)
		v, err := testSource.GetValue(opt)
		defaultAssert.Nil(v)
		defaultAssert.ErrorIs(err, ErrEmptyValue)
	})

	t.Run("test regex with uncorrect envars", func(t *testing.T) {
		setup(t)
		var key, value string
		key = "TEST_FORMAT1"
		value = "test_format1"
		opt := []string{"test.format1+format2", "test-format1.format2", "test&format1"}
		createEnvVar(t, key, value)
		for _, val := range opt {
			v, err := testSource.GetValue(val)
			defaultAssert.NotEqual(v, value)
			defaultAssert.ErrorIs(err, ErrInvalidFormat)
		}
	})
}

func TestEnvSource_Name(t *testing.T) {
	t.Run("get name", func(t *testing.T) {
		setup(t)
		defaultAssert.Equal(testSource.Name(), "ENV")
	})
}
