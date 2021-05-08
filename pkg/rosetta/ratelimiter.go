package rosetta

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/zekroTJA/timedmap"
)

var ErrRateLimited = errors.New("rate limited")

// Bucket implements a simple token bucket rate limiter.
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

// Limiter provides a limiter instance.
type Limiter interface {
	// GetBucket returns a token-bucket instance from given context.
	// command can be accessed via given ctx and should also implemented LimiterConfig.
	GetBucket(ctx *Context) *Bucket
}

type internalLimiter struct {
	executions *timedmap.TimedMap
	pool       *sync.Pool
}

func newInternalLimiter(cleanupInterval time.Duration) *internalLimiter {
	return &internalLimiter{
		executions: timedmap.New(cleanupInterval),
		pool:       &sync.Pool{New: func() interface{} { return new(Bucket) }},
	}
}

func (i *internalLimiter) GetBucket(ctx *Context) *Bucket {
	key := fmt.Sprintf("%s:%s:%s", ctx.Command.Name, ctx.Event.Author.ID, ctx.Event.Message.GuildID)
	expired := time.Duration(ctx.Command.GetLimiterBurst()) * ctx.Command.GetLimiterRestoration()

	limiter, ok := i.executions.GetValue(key).(*Bucket)
	if ok {
		if err := i.executions.SetExpire(key, expired); err != nil {
			panic(err)
		}
		return limiter
	}

	limiter = i.pool.Get().(*Bucket).setParams(ctx.Command.GetLimiterBurst(), ctx.Command.GetLimiterRestoration())
	i.executions.Set(key, limiter, expired, func(val interface{}) { i.pool.Put(val) })
	return limiter
}

// RateLimiter holds rate limiter instance.
// This Limiter will manage all instance of given rate limiter as middleware.
type RateLimiter struct {
	limiter Limiter
}

// NewRateLimiter returns a instance of RateLimiter.
//
// Additionally, one can pass a custom Limiter instance if you want to handle executions differently
// than the standard implementation.
// WIP: implements a chain of limiter.
func NewRateLimiter(limiters ...Limiter) *RateLimiter {
	var l Limiter

	if len(limiters) > 0 && limiters[0] != nil {
		l = limiters[0]
	} else {
		l = newInternalLimiter(10 * time.Minute)
	}

	return &RateLimiter{limiter: l}
}

// HandleExecution will check if the execution of ctx is valid, then send cmd downstream.
// Returns true if so, false and error otherwise.
func (r *RateLimiter) HandleExecution(ctx *Context) (bool, error) {
	limiter := r.limiter.GetBucket(ctx)
	if ok, next := limiter.Take(); !ok {
		err := ctx.RespondTextEmbedError(fmt.Sprintf("You are being rate limited. Please try again after %s", next.String()), "Rate Limiter", ErrRateLimited)
		return false, err
	}
	return true, nil
}
