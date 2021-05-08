package rosetta

import (
	"testing"
	"time"
)

var TestRateLimiter = NewRateLimiter()

var testInternalLimiter = newInternalLimiter(10 * time.Minute)

func TestBucket_Take(t *testing.T) {
	const burst = 3
	const restoration = time.Second

	b := NewBucket(burst, restoration)
	for i := 0; i < burst; i++ {
		if ok, _ := b.Take(); !ok {
			t.Errorf("token %d is not retrieved", i)
		}
	}

	ok, next := b.Take()
	if ok {
		t.Error("token retrieved yet no tokens are available")
	}
	if next > restoration || next < restoration-100*time.Microsecond {
		t.Errorf("returned next value is not in error margin: %s", next)
	}
}

func TestInternalLimiter_GetBucket(t *testing.T) {
	l1 := testInternalLimiter.GetBucket(makeTestCtx())
	l2 := testInternalLimiter.GetBucket(makeTestCtx())
	l3 := testInternalLimiter.GetBucket(TestDiffCtx)

	if l1 != l2 {
		t.Errorf("l1 (%+v) != l2 (%+v)", l1, l2)
	}

	if l1 == l3 || l2 == l3 {
		t.Error("new limiter was a dups.")
	}
}

func TestRateLimiter_HandleExecution(t *testing.T) {
	m := NewRateLimiter()
	shallPass := func() {
		ok, err := m.HandleExecution(TestDiffCtx)
		if err != nil {
			t.Errorf("Rate limiter failed: %s", err.Error())
		}
		if !ok {
			t.Error("rate limiter stopped unexpectedly")
		}
	}

	for i := 0; i < TestCommand.GetLimiterBurst(); i++ {
		shallPass()
	}
}
