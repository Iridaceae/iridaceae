package rosetta

import (
	"strings"
	"time"
)

// LimiterConfig defines command that is rate limit-able.
type LimiterConfig interface {
	// GetLimiterBurst returns max amount of tokens which can be available at time.
	// Examples:
	// 		type PingCmd struct {*Command}
	//      func (p *PingCmd) GetLimiterBurst() int {
	//			return 2
	//		}
	GetLimiterBurst() int

	// GetLimiterRestoration returns duration between new token get generated.
	// Examples:
	//      func (p *PingCmd) GetLimiterRestoration() time.Duration {
	//			return 5 * time.Second
	//		}
	GetLimiterRestoration() time.Duration
}

// Command represents a single command.
type Command struct {
	LimiterConfig
	Name        string
	Aliases     []string
	Description string
	Usage       string
	Example     string
	Flags       []string
	IgnoreCase  bool
	SubCommands []*Command
	RateLimiter *RateLimiter
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

	// Prep all middleware.
	nextHandler := c.Handler
	for _, middleware := range ctx.Router.Middlewares {
		nextHandler = middleware(nextHandler)
	}

	// Run all middleware.
	nextHandler(ctx)
}
