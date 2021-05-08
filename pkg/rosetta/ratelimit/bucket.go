package ratelimit

import "time"

// Bucket implements a simple bucket containing multiple tokens.
type Bucket struct {
	burst       int
	restoration time.Duration
	tokens      int
	lastAck     time.Time
}

// NewBucket creates a new limiter with given burst and restoration value.
func NewBucket(burst int, restore time.Duration) *Bucket {
	return new(Bucket).setParams(burst, restore)
}

// Take returns true when token is available to be taken, else false as well as duration til next token is available.
func (l *Bucket) Take() (ok bool, next time.Duration) {
	tokens := l.getTokens()
	if tokens == 0 {
		next = l.restoration - time.Since(l.lastAck)
		return
	}
	l.tokens = tokens - 1
	l.lastAck = time.Now()
	ok = true
	return
}

func (l *Bucket) setParams(burst int, restore time.Duration) *Bucket {
	l.burst = burst
	l.restoration = restore
	l.tokens = burst
	l.lastAck = time.Time{}
	return l
}

func (l *Bucket) getTokens() int {
	tokens := int(time.Since(l.lastAck)/l.restoration) + l.tokens
	if tokens > l.burst {
		return l.burst
	}
	return tokens
}
