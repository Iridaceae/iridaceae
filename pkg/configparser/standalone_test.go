package configparser

import (
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

const envNameRegex = "(ENV|env)+"

var (
	testSource    Source
	testEnvSource *EnvSource
	defaultAssert *assert.Assertions

	testOptions = &Options{
		Name:        "config.options1",
		Description: "this a mock options",
		Manager:     TestParser,
	}
)

func setupConfigTest(t *testing.T) {
	t.Helper()
	defaultAssert = assert.New(t)
	addTestSource(t, &EnvSource{})
	testSource = getTestEnvSource(t)
}

func createAndRegister(t *testing.T, name, desc string, defaultValue interface{}) error {
	t.Helper()
	_, err := TestParser.Register(name, desc, defaultValue)
	TestParser.Load()
	return err
}

func addTestSource(t *testing.T, s Source) {
	t.Helper()
	TestParser.AddSource(s)
}

func getTestEnvSource(t *testing.T) Source {
	t.Helper()
	envRgx := regexp.MustCompile(envNameRegex)
	// since we want to get Env source we will check for names
	for i := len(TestParser.sources) - 1; i >= 0; i-- {
		source := TestParser.sources[i]
		if val := envRgx.MatchString(source.Name()); val {
			return source
		}
	}

	t.Error("cannot find envsource in given manager")
	return testEnvSource
}

func createTestEnvVars(t *testing.T, key, value string) {
	t.Helper()
	err := os.Setenv(key, value)
	if err != nil {
		t.Errorf("error creating envars %s: %s", key, err.Error())
	}
}
