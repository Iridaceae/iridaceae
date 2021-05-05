package jog

import (
	"time"

	"github.com/zekroTJA/timedmap"
)

// RateLimiter represents a general rate limiter.
type RateLimiter interface {
	NotifyExecution(ctx *Context) bool
}

// StdRateLimiter represents an internal rate limiter.
type StdRateLimiter struct {
	Cooldown           time.Duration
	RateLimitedHandler ExecutionHandler
	executions         *timedmap.TimedMap
}

// NewRateLimiter creates a default rate limiter.
func NewRateLimiter(cooldown, cleanupInterval time.Duration, onRateLimited ExecutionHandler) RateLimiter {
	return &StdRateLimiter{
		Cooldown:           cooldown,
		RateLimitedHandler: onRateLimited,
		executions:         timedmap.New(cleanupInterval),
	}
}

// NotifyExecution notifies the rate limiter about a new execution and returns whether or not the execution is allowed.
func (rl *StdRateLimiter) NotifyExecution(ctx *Context) bool {
	if rl.executions.Contains(ctx.Event.Author.ID) {
		if rl.RateLimitedHandler != nil {
			next, err := rl.executions.GetExpires(ctx.Event.Author.ID)
			if err != nil {
				ctx.ObjectsMap.Set("jog_nextExecution", next)
			}
			rl.RateLimitedHandler(ctx)
		}
		return false
	}
	rl.executions.Set(ctx.Event.Author.ID, time.Now().UnixNano()/1e6, rl.Cooldown)
	return true
}
