// Package ratelimit provides a basic token-bucket limiter for rosetta router.
// This can be used to prevent users from spamming the bot and overload discord API.
package ratelimit

import (
	"time"

	"github.com/Iridaceae/iridaceae/pkg/rosetta"
)

// RateLimiter implements a managers of rate limiters.
// This can also be parsed a custom Manager instance if you want to handle limiters differently.
type RateLimiter struct {
	m       Manager
	handler rosetta.ExecutionHandler
}

// New returns a new instance of Rate Limiter.
func New(onRateLimited rosetta.ExecutionHandler, m ...Manager) *RateLimiter {
	var man Manager
	if len(m) > 0 && m[0] != nil {
		man = m[0]
	} else {
		man = newManager(10 * time.Minute)
	}
	return &RateLimiter{m: man, handler: onRateLimited}
}

func (r *RateLimiter) Handle(ctx *rosetta.Context) (bool, error) {
	limiter := r.m.GetBucket(ctx)
	if ok, next := limiter.Take(); !ok {
		r.handler(ctx, next)
		return false, rosetta.ErrRateLimited
	}
	return true, nil
}

func (r *RateLimiter) GetLayer() rosetta.MiddlewareLayer {
	return rosetta.LayerBeforeCommand
}
