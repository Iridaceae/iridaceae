package ratelimit

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	helpers2 "github.com/Iridaceae/iridaceae/internal/helpers"

	"github.com/Iridaceae/iridaceae/pkg/rosetta"

	"github.com/stretchr/testify/assert"
)

const count = 1000

var tm *managerImpl

func init() {
	_ = helpers2.LoadGlobalEnv()
	tm = newInternalManager(10 * time.Minute)
}

func loop(cmd TestCmd, f func(i int, key string, b *Bucket)) {
	wg := new(sync.WaitGroup)
	wg.Add(count)
	for i := 0; i < count; i++ {
		key := fmt.Sprintf("%s:user%d:gid%d", cmd.GetDomain(), i, i)
		b := tm.pool.Get().(*Bucket).setParams(cmd.GetLimiterBurst(), cmd.GetLimiterRestoration())
		go func(i int, key string, b *Bucket) {
			defer wg.Done()
			f(i, key, b)
		}(i, key, b)
	}
	wg.Wait()
}

func setupBucket() {
	cmd := TestCmd{false, false, false}
	loop(cmd, func(i int, key string, b *Bucket) {
		tm.executions.Set(key, b, time.Duration(cmd.GetLimiterBurst())*cmd.GetLimiterRestoration(), func(val interface{}) { tm.pool.Put(val) })
	})
}

func TestManager_GetExecutions(t *testing.T) {
	M := newInternalManager(10 * time.Minute)
	assert.Equal(t, M.executions, M.GetExecutions())
}

func TestManager_GetBucket(t *testing.T) {
	t.Run("test for dups", func(t *testing.T) {
		m := newInternalManager(10 * time.Minute)

		test := &TestCmd{false, false, false}
		l1 := m.GetBucket(test, "u1", "guild")
		l2 := m.GetBucket(test, "u1", "guild")
		l3 := m.GetBucket(test, "u2", "guild")
		assert.Equal(t, l1, l2)
		if l3 == l1 || l3 == l2 {
			t.Errorf("%+v was a duplicate of %+v & %+v.", l3, l1, l2)
		}
	})
	t.Run("add unknown bucket", func(t *testing.T) {
		setupBucket()
		cmd := &TestCmd{false, false, false}
		key := fmt.Sprintf("%s:user1:gid1", cmd.GetDomain())
		assert.Equal(t, tm.executions.GetValue(key).(*Bucket), tm.GetBucket(cmd, "user1", "gid1"))
		assert.Equal(t, count, tm.executions.Size())
	})
}

type TestCmd struct {
	wasExecuted bool
	fail        bool
	isGlobal    bool
}

func (t *TestCmd) GetLimiterBurst() int {
	return 3
}

func (t *TestCmd) GetLimiterRestoration() time.Duration {
	return time.Second
}

func (t *TestCmd) IsLimiterGlobal() bool {
	return t.isGlobal
}

func (t *TestCmd) GetInvokers() []string {
	return []string{"ping", "p"}
}

func (t *TestCmd) GetDescription() string {
	return "ping pong hello world"
}

func (t *TestCmd) GetUsage() string {
	return "`ping` - ping"
}

func (t *TestCmd) GetGroup() string {
	return rosetta.GroupFun
}

func (t *TestCmd) GetDomain() string {
	return "test.fun.ping"
}

func (t *TestCmd) GetSubPermissionRules() []rosetta.SubPermission {
	return nil
}

func (t *TestCmd) IsExecutableInDM() bool {
	return true
}

func (t *TestCmd) Exec(_ rosetta.Context) error {
	t.wasExecuted = true
	if t.fail {
		return errors.New("test error")
	}
	return nil
}
