package rosetta

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func setSubCmd(t *testing.T) {
	t.Helper()
	TestCommand.SubCommands = []*Command{{
		Name:        "t1",
		Aliases:     []string{"t1"},
		Description: "subcmd 1",
		Usage:       "t1 something_here",
		Example:     "testcmd 1 t1",
		IgnoreCase:  false,
	}}
}

func TestCommand_GetSubCmd(t *testing.T) {
	t.Run("get nil subcmd", func(t *testing.T) {
		tcmd := TestCommand.GetSubCmd("nothing")
		assert.Nil(t, tcmd)
		assert.Equal(t, 3*time.Second, tcmd.GetLimiterRestoration())
	})

	t.Run("get a subcmd", func(t *testing.T) {
		setSubCmd(t)
		tcmd := TestCommand.GetSubCmd("t1")
		assert.NotNil(t, tcmd)
		assert.Equal(t, tcmd.Name, "t1")
		assert.Equal(t, 3, tcmd.GetLimiterBurst())
	})
}
