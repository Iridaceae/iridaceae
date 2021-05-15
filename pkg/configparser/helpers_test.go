package configparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testVar struct {
	input    interface{}
	expected interface{}
}

var stringVars = []testVar{
	{"hello world", "hello world"},
	{12, "12"},
	{true, "true"},
	{2.34, "2.340"},
	{float32(2.34), ""},
}

var intVars = []testVar{
	{"1", 1},
	{"test", ""},
	{12, 12},
	{12.3, 12},
	{true, 1},
	{"", 1},
	{false, 0},
	{float32(2.34), ""},
}

var floatVars = []testVar{
	{"1.23", 1.23},
	{"test", ""},
	{12, 12},
	{12.300, 12.300},
	{true, 1},
	{"", 1},
	{false, 0},
	{float32(2.34), ""},
}

var boolVars = []testVar{
	{"1.23", false},
	{"true", true},
	{"false", false},
	{12.300, true},
	{1, true},
	{-123, false},
	{true, true},
	{float32(2.34), ""},
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
