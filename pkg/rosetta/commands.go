package rosetta

import (
	"strings"
	"time"
)

// LimitedConfig defines command that is rate limit-able.
type LimitedConfig interface {

	// GetLimiterBurst returns max amount of tokens which can be available at time.
	GetLimiterBurst() int

	// GetLimiterRestoration returns duration between new token get generated.
	GetLimiterRestoration() time.Duration
}

// Command represents a single command for given context. Command are by default implements LimitedConfig.
type Command struct {
	LimitedConfig
	Name        string
	Aliases     []string
	Description string
	Usage       string
	Example     string
	Flags       []string
	IgnoreCase  bool
	SubCommands []*Command
	Handler     ExecutionHandler
}

func (c *Command) GetLimiterBurst() int {
	return 3
}

func (c *Command) GetLimiterRestoration() time.Duration {
	return 3 * time.Second
}

// GetSubCmd returns sub command of given name if exists, else nil.
func (c *Command) GetSubCmd(name string) *Command {
	for _, sub := range c.SubCommands {
		toCheck := getIdentifiers(sub)

		// Check prefix of given string.
		if arrayContains(toCheck, name, sub.IgnoreCase) {
			return sub
		}
	}
	return nil
}

func (c *Command) trigger(ctx *Context) {
	if len(ctx.Arguments.args) > 0 {
		argument := ctx.Arguments.Get(0).Raw()
		sub := c.GetSubCmd(argument)
		if sub != nil {
			// Define arg for sub commands
			arguments := ParseArguments("")
			if ctx.Arguments.Len() > 1 {
				arguments = ParseArguments(strings.Join(strings.Split(ctx.Arguments.Raw(), " ")[1:], " "))
			}

			// Trigger subcommands
			sub.trigger(&Context{
				Session:    ctx.Session,
				Event:      ctx.Event,
				Arguments:  arguments,
				ObjectsMap: ctx.ObjectsMap,
				Router:     ctx.Router,
				Command:    sub,
			})
			return
		}
	}
	// we will wrap middleware execution before and after then command handler itself.
	if !ctx.Router.executeMiddlewares(ctx, LayerBeforeCommand) {
		return
	}
	// execute commands handlers here.
	c.Handler(ctx)
	if !ctx.Router.executeMiddlewares(ctx, LayerAfterCommand) {
		return
	}
}
