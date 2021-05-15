package configmanager

import (
	"fmt"
)

const OptionsRegex string = "^(([\\w\\.])+(\\.)([\\w]){2,4}([\\w]*))*$"

var (
	ErrEmptyValue          = fmt.Errorf("empty strings")
	ErrInvalidFormat       = fmt.Errorf("invalid format")
	ErrInvalidOptionsMatch = fmt.Errorf("invalid options match")
)

// Source acts as a generic type for different source of configs.
type Source interface {

	// GetValue will return our requested value from key.
	// This will throw error if none was found or given key doesn't
	// follow our regex parsing.
	GetValue(key string) (interface{}, error)

	// Name returns name of given source.
	Name() string
}

// Manager acts as a generic interface for us to manage config in various services.
type Manager interface {

	// RegisterSource allows users to append given configparser source to the manager.
	RegisterSource(s Source)

	// RegisterOption will add given configs to the general manager.
	RegisterOption(name, desc string, defaultValue interface{}) (Options, error)

	// LoadOptions will configure our options value into given manager.
	LoadOptions()

	// Clear will reset our sources and options mapping.
	Clear(source bool, options bool)
}

// Options holds a way for us to interact with config variables and provide
// a structured way to work with parsed envars or variable from yaml.
type Options interface {

	// LoadValue will add a given value into our manager instance.
	LoadValue()

	// UpdateValue updates a value in our current instance.
	UpdateValue(val interface{})

	// GetValue returns loaded value of given options.
	GetValue() interface{}

	// GetName returns name of given options
	GetName() string

	// ToString returns string representation of the loaded value.
	ToString() string

	// ToInt returns int representation of the loaded value.
	ToInt() int

	// ToBool returns bool representation of the loaded value.
	ToBool() bool

	// ToFloat returns float64 representation of the loaded value.
	ToFloat() float64
}
