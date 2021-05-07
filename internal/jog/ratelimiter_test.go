package jog

import "time"

var TestRateLimiter = NewRateLimiter(5*time.Second, 1*time.Second, func(ctx *Context) {
	if err := ctx.RespondText("Rate limited! try again after 5s"); err != nil {
		return
	}
})
