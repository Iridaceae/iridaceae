package rosetta

import (
	"errors"
	"testing"
	"time"

	"github.com/Iridaceae/iridaceae/pkg/util"

	"github.com/Iridaceae/iridaceae/pkg"

	"github.com/stretchr/testify/assert"
)

var (
	ErrNothingWrong = errors.New("test error")
	ctx             *Context
)

func init() {
	_ = pkg.LoadGlobalEnv()
	ctx = makeTestCtx()
}

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

func TestCommand_Trigger(t *testing.T) {
	TestCommand.trigger(ctx)
}

func TestContext_RespondText(t *testing.T) {
	ctx.Session = util.MakeTestSession()
	err := ctx.RespondText("hello world")
	assert.Nil(t, err)
}

func TestContext_RespondEmbed(t *testing.T) {
	ctx.Session = util.MakeTestSession()
	err := ctx.RespondEmbed(TestEmbedMsg)
	assert.Nil(t, err)
}

func TestContext_RespondEmbedError(t *testing.T) {
	ctx.Session = util.MakeTestSession()
	err := ctx.RespondEmbedError(ErrRateLimited)
	assert.Nil(t, err)
}

func TestContext_RespondTextEmbed(t *testing.T) {
	ctx.Session = util.MakeTestSession()
	err := ctx.RespondTextEmbed("Hello, this is a test with text and embed", TestEmbedMsg)
	assert.Nil(t, err)
}

func TestContext_RespondTextEmbedError(t *testing.T) {
	ctx.Session = util.MakeTestSession()
	err := ctx.RespondTextEmbedError("Hello, this is a error text embed test", "error response", ErrNothingWrong)
	assert.Nil(t, err)
}
