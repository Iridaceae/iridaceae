package configparser

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestParser acts as a test config manager.
var TestParser = NewManager()

var (
	TestEnvSource *EnvSource
	EnvNameRegex  = "(ENV|env)+"
	defaultAssert *assert.Assertions
	testSource    Source
)

var testOptions = &Options{
	Name:         "configtest.options1",
	Description:  "this a mock options",
	DefaultValue: nil,
	Manager:      TestParser,
}

func setup(t *testing.T) {
	t.Helper()
	defaultAssert = assert.New(t)
	addSource(t, &EnvSource{})
	testSource = getEnvSource(t)
}

func createAndRegister(t *testing.T, name, desc string, defaultValue interface{}) error {
	t.Helper()
	_, err := TestParser.Register(name, desc, defaultValue)
	TestParser.Load()
	return err
}

func addSource(t *testing.T, s Source) {
	t.Helper()
	TestParser.AddSource(s)
}

func getEnvSource(t *testing.T) Source {
	t.Helper()
	// since we want to get Env source we will check for names
	for i := len(TestParser.sources) - 1; i >= 0; i-- {
		msource := TestParser.sources[i]
		if val, _ := regexp.MatchString(EnvNameRegex, msource.Name()); val {
			return msource
		}
	}

	t.Error("cannot find envsource in given manager")
	return TestEnvSource
}

func createEnvVar(t *testing.T, key, value string) {
	t.Helper()
	err := os.Setenv(key, value)
	if err != nil {
		t.Errorf("error creating envars %s: %s", key, err.Error())
	}
}

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	t.Helper()
	assert.Equal(t, fmt.Sprintf("%+v", a), fmt.Sprintf("%+v", b))
}

func assertErr(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Errorf("didn't get error when error is expected")
	}
}
