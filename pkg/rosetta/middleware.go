package rosetta

// MiddlewareLayer defines layer in which middleware should live and executed during
// command parsing and handling.
type MiddlewareLayer int

const (
	LayerBeforeCommand MiddlewareLayer = 1 << iota
	LayerAfterCommand
)

type Middleware interface {

	// Handle is called before execution of a command handler
	// and will be passed context instance.
	//
	// if return bool is false then command handler shall not execute.
	//
	// An error should only be returned when middleware handler failed unexpectedly.
	Handle(cmd Command, ctx Context, layer MiddlewareLayer) (bool, error)

	// GetLayer returns the layer in which the middleware live.
	//
	// can be defined as bitmask, thus one can combine multiple layers to execute
	// the middleware at different point.
	GetLayer() MiddlewareLayer
}
