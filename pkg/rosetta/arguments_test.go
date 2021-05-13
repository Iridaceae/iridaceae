package rosetta

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var TestArgument = &Arguments{
	raw:  "test string",
	args: []Argument{"str1", "str2"},
}

func TestArgumentType(t *testing.T) {
	tests := []struct {
		name   string
		arg    Argument
		asFunc func(a Argument) error
	}{
		{"not bool", Argument("not_bool"), func(a Argument) error {
			_, err := a.AsBool()
			return err
		}},
		{"bool", Argument("TRUE"), func(a Argument) error {
			_, err := a.AsBool()
			return err
		}},
		{"not int", Argument("TRUE"), func(a Argument) error {
			_, err := a.AsInt()
			return err
		}},
		{"int", Argument("12"), func(a Argument) error {
			_, err := a.AsInt()
			return err
		}},
		{"not int64", Argument("TRUE"), func(a Argument) error {
			_, err := a.AsInt64()
			return err
		}},
		{"int64", Argument("1212347861234"), func(a Argument) error {
			_, err := a.AsInt64()
			return err
		}},
		{"not duration", Argument("TRUE"), func(a Argument) error {
			_, err := a.AsDuration()
			return err
		}},
		{"duration", Argument("1h30m5s"), func(a Argument) error {
			_, err := a.AsDuration()
			return err
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.asFunc(tt.arg)
			if strings.Contains(tt.name, "not") {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestArgument_Raw(t *testing.T) {
	t.Run("parse raw strings", func(t *testing.T) {
		test := Argument("hello world")
		expected := "hello world"
		assert.Equal(t, expected, test.String())
	})
}

func TestArgument_AsChannelMentionID(t *testing.T) {
	t.Run("invalid channel mention id", func(t *testing.T) {
		cid := Argument("#asdf")
		assert.Equal(t, "", cid.AsChannelMentionID())
	})
	t.Run("valid channel mention id", func(t *testing.T) {
		cid := Argument("<#128394761>")
		assert.Equal(t, "128394761", cid.AsChannelMentionID())
	})
}

func TestArgument_AsRoleMentionID(t *testing.T) {
	t.Run("invalid role mention id", func(t *testing.T) {
		cid := Argument("#asdf")
		assert.Equal(t, cid.AsRoleMentionID(), "")
	})
	t.Run("valid role mention id", func(t *testing.T) {
		cid := Argument("<@&12839476123478>")
		assert.Equal(t, "12839476123478", cid.AsRoleMentionID())
	})
}

func TestArgument_AsUserMentionID(t *testing.T) {
	t.Run("invalid user mention id", func(t *testing.T) {
		cid := Argument("#asdf")
		assert.Equal(t, cid.AsUserMentionID(), "")
	})
	t.Run("valid user mention id", func(t *testing.T) {
		cid := Argument("<@!12839476123478>")
		assert.Equal(t, "12839476123478", cid.AsUserMentionID())
	})
}

func TestParseArguments(t *testing.T) {
	t.Run("parse normal args", func(t *testing.T) {
		msg := "This is a normal message that will be parsed separated by space"
		testArguments := ParseArguments(msg)
		assert.Equal(t, msg, testArguments.Raw())
		assert.Equal(t, len(ArgumentsRegex.FindAllString(msg, -1)), testArguments.Len())
	})
	t.Run("parse normal with prefix", func(t *testing.T) {
		msg := "This is a `normal` message -that will-be parsed 'separated-by-space'"
		testArguments := ParseArguments(msg)
		assert.Equal(t, msg, testArguments.Raw())
		assert.Equal(t, len(ArgumentsRegex.FindAllString(msg, -1)), testArguments.Len())
	})
}

func TestArguments_Get(t *testing.T) {
	t.Run("get an arg", func(t *testing.T) {
		msg := "This is a normal message that will be parsed separated by space"
		testArguments := ParseArguments(msg)
		assert.Equal(t, Argument("This"), testArguments.Get(0))
	})
	t.Run("get an invalid arguments", func(t *testing.T) {
		msg := "This is a normal message that will be parsed separated by space"
		testArguments := ParseArguments(msg)
		assert.Equal(t, Argument(""), testArguments.Get(-12))
	})
	t.Run("get raw arguments", func(t *testing.T) {
		msg := "a b c"
		testargs := ParseArguments(msg)
		assert.Equal(t, "a b c", testargs.Raw())
	})
}

func TestArguments_Remove(t *testing.T) {
	t.Run("remove an arg", func(t *testing.T) {
		msg := "This is a normal message that will be parsed separated by space"
		testArguments := ParseArguments(msg)
		testArguments.Remove(0)
		assert.NotEqual(t, Argument("This"), testArguments.Get(0))
	})
	t.Run("remove a number longer than message", func(t *testing.T) {
		msg := "This is a normal message that will be parsed separated by space"
		testArguments := ParseArguments(msg)
		testArguments.Remove(100)
		assert.Equal(t, Argument("This"), testArguments.Get(0))
	})
}

func TestArguments_AsCodeblock(t *testing.T) {
	t.Run("check message as a codeblock", func(t *testing.T) {
		msg := "s```python\nprint('Hello World')```"
		testArguments := ParseArguments(msg)
		cb := testArguments.AsCodeblock()
		assert.Equal(t, "python", cb.Language)
	})
	t.Run("valid message as a codeblock without a language", func(t *testing.T) {
		msg := "s```print('Hello World')```"
		testArguments := ParseArguments(msg)
		cb := testArguments.AsCodeblock()
		assert.Equal(t, "", cb.Language)
	})
	t.Run("inline message", func(t *testing.T) {
		msg := "s`print('Hello World')`"
		testArguments := ParseArguments(msg)
		cb := testArguments.AsCodeblock()
		assert.Equal(t, "", cb.Language)
		assert.Equal(t, "print('Hello World')", cb.Content)
	})
	t.Run("invalid codeblock", func(t *testing.T) {
		msg := ""
		testArguments := ParseArguments(msg)
		cb := testArguments.AsCodeblock()
		assert.Nil(t, cb)
	})
}

func TestArguments_AsSingle(t *testing.T) {
	msg := TestArgument.AsSingle()
	assert.Equal(t, TestArgument.raw, msg.raw)
	assert.Nil(t, msg.args)
}

func TestArguments_Len(t *testing.T) {
	lenArgs := TestArgument.Len()
	assert.Equal(t, 2, lenArgs)
}

func TestArguments_Raw(t *testing.T) {
	test := Argument("test")
	assert.Equal(t, "test", test.String())
}

func TestArguments_IndexOf(t *testing.T) {
	t.Run("return nil index", func(t *testing.T) {
		notExists := TestArgument.IndexOf("this doesn't exist")
		assert.Negative(t, notExists)
	})
	t.Run("exists index", func(t *testing.T) {
		exists := TestArgument.IndexOf("str1")
		assert.Equal(t, 0, exists)
	})
}
