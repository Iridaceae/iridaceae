package configparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewManager(t *testing.T) {
	t.Run("Empty manager", func(t *testing.T) {
		assert.Equal(t, len(NewConfigManager().Options), 0)
	})
}

func TestManager_Load(t *testing.T) {
	t.Run("load a mock options", func(t *testing.T) {
		TestParser.Load()
		// should be equal to zero since we haven't register any testOptions
		assert.Equal(t, len(TestParser.Options), 0)
	})
}

func TestRegister(t *testing.T) {
	t.Run("register an unvalid options to default config manager", func(t *testing.T) {
		opt, err := Register("test-asdf", "this shouldn't register", nil)
		assert.Error(t, err)
		assert.Nil(t, opt)
	})
}

func TestLoad(t *testing.T) {
	t.Run("mock load", func(t *testing.T) {
		// we didn't actually have any config loaded so len(options) = 0
		Load()
		assert.Equal(t, len(Standalone.Options), 0)
	})
}

func TestAddSource(t *testing.T) {
	t.Run("add envsources", func(t *testing.T) {
		AddSource(&EnvSource{})
		assert.Equal(t, len(Standalone.sources), 1)
	})
}

func TestManager_AddSource(t *testing.T) {
	// NOTE: for future reference, when add more sources such as redispool and kubernetes add more test case here.
	t.Run("add envsources", func(t *testing.T) {
		m := NewConfigManager()
		m.AddSource(&EnvSource{})
		assert.Equal(t, len(m.sources), 1)
	})
}

func TestManager_Register(t *testing.T) {
	t.Run("manager with one option and non nil default value", func(t *testing.T) {
		err := createAndRegister(t, "configparser.test.options", "This is a configparser test options, that is parsed directly to TestParser", "test configs")
		assert.Equal(t, len(TestParser.Options), 1)
		assert.Nil(t, err)
	})
	t.Run("manager with one invalid option", func(t *testing.T) {
		err := createAndRegister(t, "configparser-test.options", "This shouldn't passed", nil)
		// this is options from previous test
		assert.Equal(t, len(TestParser.Options), 1)
		assert.Error(t, err)
	})
	t.Run("manager with two valid option", func(t *testing.T) {
		err := createAndRegister(t, "configparser.test.options2", "This should also parsed", nil)
		// this is options from previous test
		assert.Equal(t, len(TestParser.Options), 2)
		assert.Nil(t, err)
	})
}

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
		testOptions.LoadValue()
		assert.Nil(t, testOptions.LoadedValue)
	})

	t.Run("Load a memory value", func(t *testing.T) {
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

	t.Run("load default bool vaule", func(t *testing.T) {
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
	testOptions.LoadValue()
	t.Run("check default value", func(t *testing.T) {
		assert.Nil(t, testOptions.DefaultValue)
	})
	t.Run("update loaded value to boolean", func(t *testing.T) {
		testOptions.UpdateValue(true)
		assert.Equal(t, testOptions.LoadedValue, true)
	})
	t.Run("update loaded value to string", func(t *testing.T) {
		tstStr := "teststring"
		testOptions.UpdateValue(tstStr)
		assert.Equal(t, testOptions.LoadedValue, tstStr)
	})
	t.Run("update loaded value to int", func(t *testing.T) {
		tstInt := 1342
		testOptions.UpdateValue(tstInt)
		assert.Equal(t, testOptions.LoadedValue, tstInt)
	})
	t.Run("update loaded value to float64", func(t *testing.T) {
		tstInt := float64(12)
		testOptions.UpdateValue(tstInt)
		assert.Equal(t, testOptions.LoadedValue, tstInt)
	})
}
