package jog

import (
	"fmt"
	"sync/atomic"
	"time"
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
