package configmanager

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
		mockOption.LoadValue()
		assert.Nil(t, mockOption.LoadedValue)
	})

	t.Run("Load a memory value", func(t *testing.T) {
		TestParser.Reset()
		opt := Options{
			Name:         "random.test",
			Description:  "0x48",
			DefaultValue: 0x48,
			Manager:      TestParser,
		}
		opt.Manager.AddSource(&EnvSource{})

		opt.LoadValue()
		assert.Equal(t, opt.LoadedValue, 0x48)
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
	mockOption.LoadValue()
	t.Run("check default value", func(t *testing.T) {
		assert.Nil(t, mockOption.DefaultValue)
	})
	t.Run("update loaded value to boolean", func(t *testing.T) {
		mockOption.UpdateValue(true)
		assert.Equal(t, mockOption.LoadedValue, true)
	})
	t.Run("update loaded value to string", func(t *testing.T) {
		tstStr := "test"
		mockOption.UpdateValue(tstStr)
		assert.Equal(t, mockOption.LoadedValue, tstStr)
	})
	t.Run("update loaded value to int", func(t *testing.T) {
		tstInt := 1342
		mockOption.UpdateValue(tstInt)
		assert.Equal(t, mockOption.LoadedValue, tstInt)
	})
	t.Run("update loaded value to float64", func(t *testing.T) {
		tstInt := float64(12)
		mockOption.UpdateValue(tstInt)
		assert.Equal(t, mockOption.LoadedValue, tstInt)
	})
}

type testVar struct {
	input    interface{}
	expected interface{}
}

type testString string

func (t testString) String() string {
	return string(t)
}

var stringVars = []testVar{
	{"hello world", "hello world"},
	{testString("hello world"), "hello world"},
	{12, "12"},
	{true, "true"},
	{2.34, "2.340"},
	{float32(2.34), ""},
}

var intVars = []testVar{
	{"1", 1},
	{"test", 0},
	{12, 12},
	{12.3, 12},
	{true, 1},
	{false, 0},
	{&testVar{}, ""},
}

var floatVars = []testVar{
	{"1.23", 1.23},
	{12, float64(12)},
	{12.300, 12.300},
	{true, 1},
	{"", 0},
	{false, 0},
	{&testVar{}, ""},
}

var boolVars = []testVar{
	{"1.23", false},
	{"true", true},
	{"false", false},
	{12.300, true},
	{1, true},
	{-123, false},
	{true, true},
	{&testVar{}, false},
}

func TestConvertVal(t *testing.T) {
	impls := []struct {
		name    string
		testArr []testVar
		f       func(interface{}) interface{}
	}{
		{"convert to string", stringVars, func(i interface{}) interface{} { return toStrVal(i) }},
		{"convert to int", intVars, func(i interface{}) interface{} { return toIntVal(i) }},
		{"convert to float64", floatVars, func(i interface{}) interface{} { return toFloat64Val(i) }},
		{"convert to bool", boolVars, func(i interface{}) interface{} { return toBoolVal(i) }},
	}
	for _, tt := range impls {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Error("it didn't panic at all")
				}
			}()
			for _, i := range tt.testArr {
				assert.Equal(t, i.expected, tt.f(i.input))
			}
		})
	}
}

func TestMatchOptionsRegex(t *testing.T) {
	t.Run("key follows options", func(t *testing.T) {
		keys := "test.options.anotherone"
		ok, err := matchOptionsRegex(keys)
		assert.True(t, ok)
		assert.Nil(t, err)
	})
	t.Run("key does not options", func(t *testing.T) {
		keys := "test-options-anotherone"
		ok, err := matchOptionsRegex(keys)
		assert.False(t, ok)
		assert.ErrorIs(t, err, ErrInvalidOptionsMatch)
	})
}
