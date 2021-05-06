package jog

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	testHandlerCalled = int32(0)
	TestRateLimiter   = NewRateLimiter(5*time.Second, 1*time.Second, func(ctx *Context) {
		if err := ctx.RespondText("Rate limited! try again after 5s"); err != nil {
			return
		}
	})
	TestCommand = &Command{
		Name:        "testcmd",
		Aliases:     []string{"tcmd", "testcmd"},
		Description: "this is a test command that will increment by one then send response to context",
		Usage:       "testcmd int",
		Example:     "testcmd 1",
		IgnoreCase:  true,
		RateLimiter: TestRateLimiter,
		Handler:     testCommand,
	}
)

func testCommand(ctx *Context) {
	atomic.AddInt32(&testHandlerCalled, 1)
	if err := ctx.RespondText(fmt.Sprintf("handlercalled value: %d", testHandlerCalled)); err != nil {
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
