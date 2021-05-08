package ratelimit

import (
	"fmt"
	"sync"
	"time"

	"github.com/zekroTJA/timedmap"

	"github.com/Iridaceae/iridaceae/pkg/rosetta"
)

// Manager holds multiple buckets, act as an interface for rate limiter.
type Manager interface {

	// GetBucket returns a token-bucket rate limiter from given context.
	GetBucket(ctx *rosetta.Context) *Bucket
}

type manager struct {
	executions *timedmap.TimedMap
	pool       *sync.Pool
}

func newManager(cleanupInterval time.Duration) *manager {
	return &manager{
		executions: timedmap.New(cleanupInterval),
		pool:       &sync.Pool{New: func() interface{} { return new(Bucket) }},
	}
}

func (m *manager) GetBucket(ctx *rosetta.Context) *Bucket {
	key := fmt.Sprintf("%s:%s:%s", ctx.Command.Name, ctx.Event.Author.ID, ctx.Event.Message.GuildID)
	expired := time.Duration(ctx.Command.GetLimiterBurst()) * ctx.Command.GetLimiterRestoration()

	limiter, ok := m.executions.GetValue(key).(*Bucket)
	if ok {
		_ = m.executions.SetExpire(key, expired)
		return limiter
	}

	limiter = m.pool.Get().(*Bucket).setParams(ctx.Command.GetLimiterBurst(), ctx.Command.GetLimiterRestoration())
	m.executions.Set(key, limiter, expired, func(val interface{}) { m.pool.Put(val) })
	return limiter
}
