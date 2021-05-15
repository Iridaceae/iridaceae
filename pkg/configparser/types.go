package configparser

import "fmt"

const OptionsRegex string = "^(([\\w\\.])+(\\.)([\\w]){2,4}([\\w]*))*$"

var (
	ErrEmptyValue          = fmt.Errorf("empty strings")
	ErrInvalidFormat       = fmt.Errorf("invalid format")
	ErrInvalidOptionsMatch = fmt.Errorf("invalid options match")
)

// Source acts as a generic type for different source of configs.
// ex: env, yaml, toml. refers to EnvSource for envars parsing.
// Source also have the ability to marshall and unmarshal given config to file.
type Source interface {

	// GetValue will return our requested value from key.
	// This will throw error if none was found or given key doesn't
	// follow our regex parsing.
	GetValue(key string) (interface{}, error)

	// Name returns name of given source.
	Name() string

	// Marshal serialize data stream into
}

type Manager interface {

	// AddSource allows users to append given configparser source to the manager.
	AddSource(s Source)

	// Register will add given configs to the general manager.
	Register(name, desc string, defaultValue interface{}) (*Options, error)

	// Load will configure our options value into given manager.
	Load()

	// Reset will reset our sources and options mapping.
	Reset()
}
