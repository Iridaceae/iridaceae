package jog

import (
	"github.com/Iridaceae/iridaceae/pkg"
)

var (
	TestLogger     = pkg.NewLogger(pkg.Info, "jog_testLogger")
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
