package configmanager

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestParser acts as a test configparser manager that can be used globally.
var TestParser = NewConfigManager().(*managerImpl)

var mockOption = &Options{
	Name:        "configparser.options1",
	Description: "this a mock options",
	Manager:     TestParser,
}

func setupConfigTest(t *testing.T) {
	t.Helper()
	TestParser.AddSource(&EnvSource{})
}

func createAndRegister(t *testing.T, name, desc string, defaultValue interface{}) error {
	t.Helper()
	_, err := TestParser.Register(name, desc, defaultValue)
	TestParser.Load()
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
		opt, err := Register("test-asdf", "this shouldn't register", nil)
		assert.Error(t, err)
		assert.Nil(t, opt)
	})
}

func TestLoad(t *testing.T) {
	t.Run("mock load", func(t *testing.T) {
		// we didn't actually have any configparser loaded so len(options) = 0
		Load()
		assert.Equal(t, len(Standalone.Options), 0)
	})
}

func TestAddSource(t *testing.T) {
	t.Run("add envsources", func(t *testing.T) {
		AddSource(&EnvSource{})
		assert.Equal(t, len(Standalone.sources), 1)
	})
	t.Run("reset after source", func(t *testing.T) {
		AddSource(&EnvSource{})
		Reset()
		assert.Equal(t, len(Standalone.sources), 0)
	})
}
