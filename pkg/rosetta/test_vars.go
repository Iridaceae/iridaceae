package rosetta

import (
	"errors"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/Iridaceae/iridaceae/pkg/sclog/log"
)

type TestMiddleware struct {
	name string
}

func (t *TestMiddleware) Handle(ctx *Context) (bool, error) {
	ctx.ObjectsMap.Set("myObject", 13)
	t.name = ctx.Event.ID

	// retrieve the object
	obj, ok := ctx.ObjectsMap.GetValue("myObject").(int)
	if !ok {
		return false, errors.New("null object")
	}
	log.Info("rosetta_objTest", obj)
	return true, nil
}

func (t *TestMiddleware) GetLayer() MiddlewareLayer {
	return LayerBeforeCommand
}

var (
	TestObjectsMap *ObjectsMap

	TestArgument = &Arguments{
		raw: "test string",
		args: []*Argument{
			{"str1"},
			{"str2"},
		},
	}

	TestCommand = &Command{
		Name:        "obj",
		Aliases:     []string{"object", "obj"},
		Description: "this is a test command that will inject an ObjectsMap",
		Usage:       "obj",
		Example:     "obj",
		IgnoreCase:  true,
		SubCommands: []*Command{},
		Handler:     onTestCmd,
	}

	TestEmbedMsg = &discordgo.MessageEmbed{
		Type:        "rich",
		Title:       "This is a test message",
		Description: "Embed nice",
		Timestamp:   time.Now().Format(time.RFC3339),
		Color:       0xffff00,
	}

	TestRouter = New(&Router{
		Prefixes:         []string{"!"},
		IgnorePrefixCase: true,
		BotsAllowed:      false,
		Commands:         []*Command{},
		Middlewares:      []Middleware{},
		PingHandler: func(ctx *Context, _ ...interface{}) {
			log.Info("context", ctx.Arguments)
		},
	})
)

func onTestCmd(ctx *Context, _ ...interface{}) {
}
