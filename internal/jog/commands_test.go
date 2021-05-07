package jog

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

var TestCommand = &Command{
	Name:        "obj",
	Aliases:     []string{"object"},
	Description: "this is a test command that will inject an ObjectsMap",
	Usage:       "obj",
	Example:     "obj",
	IgnoreCase:  true,
	SubCommands: []*Command{},
	RateLimiter: TestRateLimiter,
	Handler:     testCommand,
}

func testCommand(ctx *Context) {
	if err := ctx.RespondText(strconv.Itoa(ctx.ObjectsMap.GetValue("myObject").(int))); err != nil {
		return
	}
}

func TestCommand_GetSubCmd(t *testing.T) {
	t.Run("get nil subcmd", func(t *testing.T) {
		tcmd := TestCommand.GetSubCmd("nothing")
		assert.Nil(t, tcmd)
	})

	t.Run("get a subcmd", func(t *testing.T) {
		TestCommand.SubCommands = []*Command{&Command{
			Name:        "t1",
			Aliases:     []string{"t1"},
			Description: "subcmd 1",
			Usage:       "t1 something_here",
			Example:     "testcmd 1 t1",
			IgnoreCase:  false,
		}}
		tcmd := TestCommand.GetSubCmd("t1")
		assert.NotNil(t, tcmd)
		assert.Equal(t, tcmd.Name, "t1")
	})
}
