package ratelimit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	burst       int = 5
	restoration     = 3 * time.Second
)

func TestBucket_Take(t *testing.T) {
	randBucket := NewBucket(burst, restoration)
	assert.Equal(t, 5, randBucket.burst)

	for i := 0; i < burst; i++ {
		ok, _ := randBucket.Take()
		assert.True(t, ok)
	}
	ok, next := randBucket.Take()
	assert.False(t, ok)

	assert.Less(t, next, restoration)
	assert.Greater(t, next, restoration-100*time.Microsecond)
}
