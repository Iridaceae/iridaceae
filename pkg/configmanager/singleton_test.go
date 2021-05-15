package configmanager

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestParser acts as a test configparser manager that can be used globally.
var TestParser = NewDefaultManager().(*managerImpl)

var mockOption = &optionsImpl{
	Name:        "configparser.options1",
	Description: "this a mock options",
	Manager:     TestParser,
}

func setupConfigTest(t *testing.T) {
	t.Helper()
	TestParser.Clear(true, false)
	TestParser.RegisterSource(&EnvSource{})
}

func createAndRegister(t *testing.T, name, desc string, defaultValue interface{}) error {
	t.Helper()
	TestParser.Clear(true, false)
	_, err := TestParser.RegisterOption(name, desc, defaultValue)
	TestParser.LoadOptions()
	return err
}

func createTestEnvVars(t *testing.T, key, value string) {
	t.Helper()
	err := os.Setenv(key, value)
	if err != nil {
		t.Errorf("error creating envars %s: %s", key, err.Error())
	}
}

func TestRegister(t *testing.T) {
	t.Run("register an unvalid options to default configparser manager", func(t *testing.T) {
		opt, err := RegisterOption("test-asdf", "this shouldn't register", nil)
		assert.Error(t, err)
		assert.Nil(t, opt)
	})
}

func TestLoad(t *testing.T) {
	t.Run("mock load", func(t *testing.T) {
		// we didn't actually have any configparser loaded so len(options) = 0
		LoadOptions()
		assert.Equal(t, len(Standalone.Options), 0)
	})
}

func TestAddSource(t *testing.T) {
	t.Run("add envsources", func(t *testing.T) {
		RegisterSource(&EnvSource{})
		assert.Equal(t, len(Standalone.sources), 1)
	})
	t.Run("reset after source", func(t *testing.T) {
		RegisterSource(&EnvSource{})
		Clear()
		assert.Equal(t, len(Standalone.sources), 0)
	})
}
