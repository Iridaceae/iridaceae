package ratelimit

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Iridaceae/iridaceae/pkg/rosetta"
)

func TestRateLimiter_GetLayer(t *testing.T) {
	tm := New(func(ctx *rosetta.Context, _ ...interface{}) { _ = ctx.RespondText("Rate Limited") })
	assert.Equal(t, rosetta.LayerBeforeCommand, tm.GetLayer())
}

func TestRateLimiter_Handle(t *testing.T) {
	tm := New(func(ctx *rosetta.Context, _ ...interface{}) { _ = fmt.Sprintf("test ctx: %s", ctx.Command.Name) })
	pass := func() {
		ok, err := tm.Handle(TestCtx)
		assert.Nil(t, err)
		assert.True(t, ok, "rate limiter stopped unexpectedly")
	}
	fail := func() {
		ok, err := tm.Handle(TestCtx)
		assert.Error(t, err)
		assert.False(t, ok, "rate limiter passed unexpectedly")
	}

	for i := 0; i < Command.GetLimiterBurst(); i++ {
		pass()
	}
	fail()
}
