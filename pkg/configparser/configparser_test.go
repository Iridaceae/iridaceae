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
		// should be equal to zero since we haven't register any testOpts
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
