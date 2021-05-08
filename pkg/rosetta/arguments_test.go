package rosetta

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArgument_AsBool(t *testing.T) {
	t.Run("wrong return bool", func(t *testing.T) {
		test := &Argument{"non_bool"}
		_, err := test.AsBool()
		assert.Error(t, err)
	})

	t.Run("test good boolean", func(t *testing.T) {
		test := &Argument{"TRUE"}
		b, err := test.AsBool()
		assert.Nil(t, err)
		assert.Equal(t, true, b)
	})
}

func TestArgument_AsInt(t *testing.T) {
	t.Run("non int", func(t *testing.T) {
		test := &Argument{"TRUE"}
		_, err := test.AsInt()
		assert.Error(t, err)
	})

	t.Run("good int", func(t *testing.T) {
		test := &Argument{"12"}
		i, err := test.AsInt()
		assert.Nil(t, err)
		assert.Equal(t, 12, i)
	})
}

func TestArgument_AsInt64(t *testing.T) {
	t.Run("invalid int", func(t *testing.T) {
		test := &Argument{"TRUE"}
		_, err := test.AsInt()
		assert.Error(t, err)
	})

	t.Run("parse int64", func(t *testing.T) {
		test := &Argument{"1214984716"}
		i, err := test.AsInt()
		assert.Nil(t, err)
		assert.Equal(t, 1214984716, i)
	})
}

func TestArgument_AsDuration(t *testing.T) {
	t.Run("invalid duration", func(t *testing.T) {
		test := &Argument{"12384"}
		_, err := test.AsDuration()
		assert.Error(t, err)
	})

	t.Run("parse duration", func(t *testing.T) {
		test := &Argument{"1h30m5s"}
		tp, err := test.AsDuration()
		assert.Nil(t, err)
		assert.Equal(t, "1h30m5s", tp.String())
	})
}

func TestArgument_Raw(t *testing.T) {
	t.Run("parse raw strings", func(t *testing.T) {
		test := &Argument{"hello world"}
		expected := "hello world"
		assert.Equal(t, expected, test.Raw())
	})
}

func TestArgument_AsChannelMentionID(t *testing.T) {
	t.Run("invalid channel mention id", func(t *testing.T) {
		cid := &Argument{"#asdf"}
		assert.Equal(t, "", cid.AsChannelMentionID())
	})
	t.Run("valid channel mention id", func(t *testing.T) {
		cid := &Argument{"<#12839476123478>"}
		assert.Equal(t, "12839476123478", cid.AsChannelMentionID())
	})
}

func TestArgument_AsRoleMentionID(t *testing.T) {
	t.Run("invalid role mention id", func(t *testing.T) {
		cid := &Argument{"#asdf"}
		assert.Equal(t, cid.AsRoleMentionID(), "")
	})
	t.Run("valid role mention id", func(t *testing.T) {
		cid := &Argument{"<@&12839476123478>"}
		assert.Equal(t, "12839476123478", cid.AsRoleMentionID())
	})
}

func TestArgument_AsUserMentionID(t *testing.T) {
	t.Run("invalid user mention id", func(t *testing.T) {
		cid := &Argument{"#asdf"}
		assert.Equal(t, cid.AsUserMentionID(), "")
	})
	t.Run("valid user mention id", func(t *testing.T) {
		cid := &Argument{"<@!12839476123478>"}
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
		assert.Equal(t, &Argument{"This"}, testArguments.Get(0))
	})
}

func TestArguments_Remove(t *testing.T) {
	t.Run("remove an arg", func(t *testing.T) {
		msg := "This is a normal message that will be parsed separated by space"
		testArguments := ParseArguments(msg)
		testArguments.Remove(0)
		assert.NotEqual(t, &Argument{"This"}, testArguments.Get(0))
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
}
