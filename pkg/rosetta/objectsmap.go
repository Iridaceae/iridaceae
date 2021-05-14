package rosetta

// ReadOnlyObjectMap provides a thread-safe key-value map to get previous set items from.
type ReadOnlyObjectMap interface {

	// GetObject returns a value from its key.
	// Returns nil if no object is stored.
	GetObject(key string) (value interface{})

	// SetObject sets a value to the object.
	// This is used to workaround di.Container to get
	// our router instance inside our router context.
	SetObject(key string, value interface{})
}

// ObjectMap wraps around ReadOnlyObjectMap to provide a way to set value to and get value back.
type ObjectMap interface {
	ReadOnlyObjectMap
}
