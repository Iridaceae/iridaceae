package ratelimiter

import (
	"sync"
	"time"

	"github.com/TensRoses/iris/internal/configparser"

	"golang.org/x/time/rate"
)

var (
	MaxPerSeconds, _ = configparser.Register("iris.ratelimiter.maxperseconds", "Max frequency per seconds til the rate limiter starts", 3)
	MaxBurst, _      = configparser.Register("iris.ratelimiter.maxburst", "Max amounts of burst detected for rl", 3)
	// IrisRL is a generic multiratelimiter for the bot that manages different rate limiters for different tasks (cmd handlers, db, etc) and can be used right out of the box.
	//  example imports: _ "github.com/TensRoses/iris/internal/ratelimiter
	IrisRL *MultiRateLimiter
)

func init() {
	IrisRL = NewMultiRateLimiter(MaxPerSeconds.GetFloat(), MaxBurst.GetInt())
}

// MultiRateLimiter holds multiple RateLimiters depending on different tasks.
// NOTE: implements uber key-bucket rate-limiter instead of rate.Limiter.
type MultiRateLimiter struct {
	mu            sync.Mutex
	limiters      map[interface{}]*rate.Limiter
	maxPerSeconds float64
	maxBurst      int
}

// NewMultiRateLimiter defines a singleton MultiRateLimiter.
func NewMultiRateLimiter(maxPerSecond float64, maxBurst int) *MultiRateLimiter {
	return &MultiRateLimiter{
		limiters:      make(map[interface{}]*rate.Limiter),
		maxPerSeconds: maxPerSecond,
		maxBurst:      maxBurst,
	}
}

// GetLimiter returns a current RateLimiter lives inside MultiRateLimiter, create new if not found one.
func (m *MultiRateLimiter) GetLimiter(key interface{}) *rate.Limiter {
	m.mu.Lock()
	defer m.mu.Unlock()

	if current, ok := m.limiters[key]; ok {
		return current
	}

	m.limiters[key] = rate.NewLimiter(rate.Limit(m.maxPerSeconds), m.maxBurst)
	return m.limiters[key]
}

func (m *MultiRateLimiter) AllowN(key interface{}, now time.Time, n int) bool {
	return m.GetLimiter(key).AllowN(now, n)
}
