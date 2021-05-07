package jog

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Iridaceae/iridaceae/pkg"
)

var (
	TestLogger = pkg.NewLogger(pkg.Info, "jog_testLogger")
	// TestMiddleware is a middleware that will inject a object into context.
	TestMiddleware = func(next ExecutionHandler) ExecutionHandler {
		return func(ctx *Context) {
			ctx.ObjectsMap.Set("myObject", 13)

			// retrieve the object
			obj, ok := ctx.ObjectsMap.GetValue("myObject").(int)
			if !ok {
				return
			}
			pkg.StdLogger.Info("jog_objTest", obj)

			// call next execution handler
			next(ctx)
		}
	}
	TestRouter = Create(&Router{
		Prefixes:         []string{"!"},
		IgnorePrefixCase: true,
		BotsAllowed:      false,
		Logger:           TestLogger,
		Commands:         []*Command{},
		Middlewares:      []Middleware{},
		PingHandler: func(ctx *Context) {
			if err := ctx.RespondText("pong!"); err != nil {
				panic(err)
			}
		},
	})
)

func TestCreate(t *testing.T) {
	r := Create(TestRouter)
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
		t1 := TestRouter.GetCmd("test_options")
		assert.NotEqual(t, TestCommand, t1)
	})

	t.Run("register test", func(t *testing.T) {
		TestRouter.RegisterCmd(TestCommand)
		t2 := TestRouter.GetCmd("obj")
		assert.Equal(t, TestCommand, t2)
	})
}

func TestRouter_RegisterMiddleware(t *testing.T) {
	TestRouter.RegisterMiddleware(TestMiddleware)
	assert.Equal(t, 1, len(TestRouter.Middlewares))
}

func TestRouter_InitializeStorage(t *testing.T) {
	TestRouter.InitializeStorage("jog_testMain")
	assert.NotNil(t, TestRouter.Storage)
}

func TestRouter_Handler(t *testing.T) {
	ctx := makeTestCtx()
	ctx.Session = makeTestSession()

	TestRouter.RegisterCmd(TestCommand)
	TestRouter.RegisterMiddleware(TestMiddleware)
	ctx.Router = TestRouter
	TestRouter.Handler()
	assert.Equal(t, TestCommand, TestRouter.Commands[0])
}
