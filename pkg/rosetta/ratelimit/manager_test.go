package ratelimit

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/Iridaceae/iridaceae/pkg/rosetta"

	"github.com/Iridaceae/iridaceae/pkg"

	"github.com/stretchr/testify/assert"
)

func init() {
	_ = pkg.LoadGlobalEnv()
}

func TestManager_GetBucket(t *testing.T) {
	m := newManager(10 * time.Minute)

	test := &TestCmd{false, false, false}
	l1 := m.GetBucket(test, "u1", "guild")
	l2 := m.GetBucket(test, "u1", "guild")
	l3 := m.GetBucket(test, "u2", "guild")
	assert.Equal(t, l1, l2)
	if l3 == l1 || l3 == l2 {
		t.Error(errDupsBucket(l3))
	}
}

func errDupsBucket(i1 interface{}) string {
	return fmt.Sprintf("%+v was a duplicate of l1 & l2.", i1)
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

func (t *TestCmd) Exec(ctx rosetta.Context) error {
	t.wasExecuted = true
	if t.fail {
		return errors.New("test error")
	}
	return nil
}
