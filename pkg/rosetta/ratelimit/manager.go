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
	GetBucket(cmd rosetta.Command, uid, gid string) *Bucket

	// GetExecutions returns our internal execution mapping.
	GetExecutions() *timedmap.TimedMap
}

type managerImpl struct {
	executions *timedmap.TimedMap
	pool       *sync.Pool
}

func (m *managerImpl) GetExecutions() *timedmap.TimedMap {
	return m.executions
}

func newInternalManager(cleanupInterval time.Duration) *managerImpl {
	return &managerImpl{
		executions: timedmap.New(cleanupInterval),
		pool:       &sync.Pool{New: func() interface{} { return new(Bucket) }},
	}
}

func (m *managerImpl) GetBucket(cmd rosetta.Command, uid, gid string) *Bucket {
	key := fmt.Sprintf("%s:%s:%s", cmd.GetDomain(), uid, gid)

	// all command should implements LimitedConfig.
	lcmd, _ := cmd.(rosetta.LimitedConfig)
	expired := time.Duration(lcmd.GetLimiterBurst()) * lcmd.GetLimiterRestoration()

	limiter, ok := m.executions.GetValue(key).(*Bucket)
	if ok {
		_ = m.executions.SetExpire(key, expired)
		return limiter
	}

	limiter = m.pool.Get().(*Bucket).setParams(lcmd.GetLimiterBurst(), lcmd.GetLimiterRestoration())
	m.executions.Set(key, limiter, expired, func(val interface{}) { m.pool.Put(val) })
	return limiter
}
