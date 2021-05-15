package configparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptions_GetBool(t *testing.T) {
	t.Run("Get bool from an bool options", func(t *testing.T) {
		topt := &Options{
			Name:         "bool.opts",
			Description:  "boolean opts",
			DefaultValue: true,
			LoadedValue:  true,
		}
		assert.True(t, topt.GetBool())
	})
	t.Run("Get another valid bool", func(t *testing.T) {
		topt := &Options{
			Name:         "bool.opts",
			Description:  "boolean opts",
			DefaultValue: "true",
			LoadedValue:  "true",
		}
		assert.True(t, topt.GetBool())
	})
	t.Run("Get a false bool", func(t *testing.T) {
		topt := &Options{
			Name:         "bool.opts",
			Description:  "boolean opts",
			DefaultValue: "t",
			LoadedValue:  "t",
		}
		assert.False(t, topt.GetBool())
	})
	t.Run("Get a true bool from int", func(t *testing.T) {
		topt := &Options{
			Name:         "bool.opts",
			Description:  "boolean opts",
			DefaultValue: 1,
			LoadedValue:  1,
		}
		assert.True(t, topt.GetBool())
	})
}

func TestOptions_GetFloat(t *testing.T) {
	t.Run("Get a float options", func(t *testing.T) {
		fopt := &Options{
			Name:         "float.opts",
			Description:  "float opts",
			DefaultValue: 12.3,
			LoadedValue:  12.3,
		}
		assert.Equal(t, fopt.GetFloat(), 12.3)
	})
	t.Run("Get from int options", func(t *testing.T) {
		fopt := &Options{
			Name:         "float.opts",
			Description:  "float opts",
			DefaultValue: 12,
			LoadedValue:  12,
		}
		assert.Equal(t, fopt.GetFloat(), float64(12))
	})
	t.Run("Get a string options", func(t *testing.T) {
		fopt := &Options{
			Name:         "float.opts",
			Description:  "float opts",
			DefaultValue: "12.3",
			LoadedValue:  "12.3",
		}
		assert.Equal(t, fopt.GetFloat(), 12.3)
	})
}

func TestOptions_GetInt(t *testing.T) {
	t.Run("Get a int options", func(t *testing.T) {
		iopt := &Options{
			Name:         "test.int",
			Description:  "test int",
			DefaultValue: 1,
			LoadedValue:  1,
		}
		assert.Equal(t, iopt.GetInt(), 1)
	})
	t.Run("Get a int option from string", func(t *testing.T) {
		iopt := &Options{
			Name:         "test.int",
			Description:  "test int",
			DefaultValue: "1",
			LoadedValue:  "1",
		}
		assert.Equal(t, iopt.GetInt(), 1)
	})
}

func TestOptions_GetString(t *testing.T) {
	t.Run("Get a string options", func(t *testing.T) {
		sopt := &Options{
			Name:         "test.string",
			Description:  "test int",
			DefaultValue: "hello world",
			LoadedValue:  "hello world",
		}
		assert.Equal(t, sopt.GetString(), "hello world")
	})
	t.Run("Get a sttring from int", func(t *testing.T) {
		sopt := &Options{
			Name:         "test.string",
			Description:  "test int",
			DefaultValue: 12,
			LoadedValue:  12,
		}
		assert.Equal(t, sopt.GetString(), "12")
	})
}

func TestOptions_LoadValue(t *testing.T) {
	t.Run("Load a nil default value", func(t *testing.T) {
		testOpts.LoadValue()
		assert.Nil(t, testOpts.LoadedValue)
	})

	t.Run("Load a memory value", func(t *testing.T) {
		TestParser.Clear()
		opt := Options{
			Name:         "random.test",
			Description:  "0x48",
			DefaultValue: 0x48,
			Manager:      TestParser,
		}
		opt.Manager.AddSource(&EnvSource{})

		opt.LoadValue()
		assert.Equal(t, opt.LoadedValue, 0x48)
		assert.Equal(t, 1, len(opt.Manager.sources))
	})

	t.Run("load default bool value", func(t *testing.T) {
		opt := Options{
			Name:         "rand.test",
			DefaultValue: true,
			Manager:      TestParser,
		}
		opt.Manager.AddSource(&EnvSource{})
		opt.LoadValue()

		assert.Equal(t, opt.LoadedValue, true)
	})
}

func TestOptions_UpdateValue(t *testing.T) {
	testOpts.LoadValue()
	t.Run("check default value", func(t *testing.T) {
		assert.Nil(t, testOpts.DefaultValue)
	})
	t.Run("update loaded value to boolean", func(t *testing.T) {
		testOpts.UpdateValue(true)
		assert.Equal(t, testOpts.LoadedValue, true)
	})
	t.Run("update loaded value to string", func(t *testing.T) {
		tstStr := "teststring"
		testOpts.UpdateValue(tstStr)
		assert.Equal(t, testOpts.LoadedValue, tstStr)
	})
	t.Run("update loaded value to int", func(t *testing.T) {
		tstInt := 1342
		testOpts.UpdateValue(tstInt)
		assert.Equal(t, testOpts.LoadedValue, tstInt)
	})
	t.Run("update loaded value to float64", func(t *testing.T) {
		tstInt := float64(12)
		testOpts.UpdateValue(tstInt)
		assert.Equal(t, testOpts.LoadedValue, tstInt)
	})
}
