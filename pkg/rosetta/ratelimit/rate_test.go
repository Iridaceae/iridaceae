package ratelimit

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Iridaceae/iridaceae/pkg/rosetta"
)

func testLoop(t *testing.T, testCtx *rosetta.Context, m ...Manager) *RateLimiter {
	t.Helper()

	testHandler := func(ctx *rosetta.Context, _ ...interface{}) {}
	tm := New(testHandler, m...)

	pass := func() {
		ok, err := tm.Handle(testCtx)
		assert.Nil(t, err)
		assert.True(t, ok, "rate limiter stopped unexpectedly")
	}
	fail := func() {
		ok, err := tm.Handle(testCtx)
		assert.Error(t, err)
		assert.False(t, ok, "rate limiter passed unexpectedly")
	}

	for i := 0; i < Command.GetLimiterBurst(); i++ {
		pass()
	}
	fail()

	return tm
}

func TestRateLimiter_GetLayer(t *testing.T) {
	tm := testLoop(t, TestCtx)
	assert.Equal(t, rosetta.LayerBeforeCommand, tm.GetLayer())
}

// TODO: custom manager test.
func TestRateLimiter_Handle(t *testing.T) {
	t.Run("test rate limiter with default manager", func(t *testing.T) {
		_ = testLoop(t, TestCtx)
	})
}
