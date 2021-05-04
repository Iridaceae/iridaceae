package jog

import "strings"

// Command represents a single command.
type Command struct {
	Name        string
	Aliases     []string
	Description string
	Usage       string
	Example     string
	Flags       []string
	IgnoreCase  bool
	SubCommands []*Command
	RateLimiter RateLimiter
	Handler     ExecutionHandler
}

// GetSubCmd returns sub command of given name if exists, else nil.
func (c *Command) GetSubCmd(name string) *Command {
	for _, sub := range c.SubCommands {
		toCheck := make([]string, 0, len(sub.Aliases)+1)
		toCheck = append(toCheck, sub.Name)
		toCheck = append(toCheck, sub.Aliases...)

		// Check prefix of given string.
		if arrayContains(toCheck, name, sub.IgnoreCase) {
			return sub
		}
	}
	return nil
}

func (c *Command) NotifyRateLimiter(ctx *Context) bool {
	if c.RateLimiter == nil {
		return true
	}
	return c.RateLimiter.NotifyExecution(ctx)
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
