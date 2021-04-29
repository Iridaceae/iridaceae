package configparser

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestNewManager(t *testing.T) {
	t.Run("Empty manager", func(t *testing.T) {
		assertEqual(t, len(NewManager().Options), 0)
	})
}

func TestManager_AddSource(t *testing.T) {
	// NOTE: for future reference, when add more sources such as redispool and kubernetes add more test case here.
	t.Run("add envsources", func(t *testing.T) {
		m := NewManager()
		m.AddSource(&EnvSource{})
		assertEqual(t, len(m.sources), 1)
	})
}

func TestManager_Register(t *testing.T) {
	t.Run("manager with one option and non nil default value", func(t *testing.T) {
		err := createAndRegister(t, "configparser.test.options", "This is a configparser test options, that is parsed directly to TestParser", "test configs")
		assertEqual(t, len(TestParser.Options), 1)
		assert.Nil(t, err)
	})
	t.Run("manager with one invalid option", func(t *testing.T) {
		err := createAndRegister(t, "configparser-test.options", "This shouldn't passed", nil)
		// this is options from previous test
		assertEqual(t, len(TestParser.Options), 1)
		assertErr(t, err)
	})
	t.Run("manager with two valid option", func(t *testing.T) {
		err := createAndRegister(t, "configparser.test.options2", "This should also parsed", nil)
		// this is options from previous test
		assertEqual(t, len(TestParser.Options), 2)
		assert.Nil(t, err)
	})
}

func TestOptions_UpdateValue(t *testing.T) {
	testOptions.LoadValue()
	t.Run("check default value", func(t *testing.T) {
		assert.Nil(t, testOptions.DefaultValue)
	})
	t.Run("update loaded value to boolean", func(t *testing.T) {
		testOptions.UpdateValue(true)
		assertEqual(t, testOptions.LoadedValue, true)
	})
	t.Run("update loaded value to string", func(t *testing.T) {
		tstStr := "teststring"
		testOptions.UpdateValue(tstStr)
		assertEqual(t, testOptions.LoadedValue, tstStr)
	})
	t.Run("update loaded value to int", func(t *testing.T) {
		tstInt := 1342
		testOptions.UpdateValue(tstInt)
		assertEqual(t, testOptions.LoadedValue, tstInt)
	})
}
