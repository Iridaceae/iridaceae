package rosetta

import (
	"testing"

	"github.com/Iridaceae/iridaceae/pkg/util"

	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	r := New(TestRouter)
	assert.NotNil(t, r.Storage)
}

func TestRouter_RegisterCmd(t *testing.T) {
	TestRouter.RegisterCmd(TestCommand)
	assert.Equal(t, 1, len(TestRouter.Commands))
	assert.Equal(t, "obj", TestRouter.Commands[0].Name)
}

func TestRouter_GetCmd(t *testing.T) {
	t.Run("not register test", func(t *testing.T) {
		// we haven't register commands yet
		_, ok := TestRouter.GetCmd("test_options")
		assert.False(t, ok)
	})

	t.Run("register test", func(t *testing.T) {
		TestRouter.RegisterCmd(TestCommand)
		t2, _ := TestRouter.GetCmd("obj")
		assert.Equal(t, TestCommand, t2)
	})
}

func TestRouter_RegisterMiddleware(t *testing.T) {
	TestRouter.RegisterMiddleware(&TestMiddleware{})
	assert.Equal(t, 1, len(TestRouter.Middlewares))
}

func TestRouter_InitializeStorage(t *testing.T) {
	TestRouter.InitializeStorage("rosetta_testMain")
	assert.NotNil(t, TestRouter.Storage)
}

func TestRouter_Handler(t *testing.T) {
	ctx := makeTestCtx()
	ctx.Session = util.MakeTestSession()

	TestRouter.RegisterMiddleware(&TestMiddleware{})
	TestRouter.RegisterCmd(TestCommand)
	ctx.Router = TestRouter
	TestRouter.Handler()
	assert.Equal(t, TestCommand, TestRouter.Commands[0])
}

func TestRouter_RegisterDefaultHelpCommand(t *testing.T) {
	t.Run("should return no error", func(t *testing.T) {
		setSubCmd(t)
		ctx := makeTestCtx()
		s := util.MakeTestSession()
		ctx.Router = TestRouter
		TestRouter.RegisterMiddleware(&TestMiddleware{})
		TestRouter.RegisterDefaultHelpCommand(s)
		TestRouter.Initialize(s)
		assert.Equal(t, 4, len(TestRouter.Commands))

		// help command should alive
		_, ok := TestRouter.GetCmd("help")
		assert.True(t, ok)
	})
}
