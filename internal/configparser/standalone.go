package configparser

// Standalone is a singleton config manager that acts as a general manager for iris.
var Standalone = NewManager()

// AddSource wraps around Manager.AddSource and can be accessed via configparser.AddSource.
func AddSource(s Source) {
	Standalone.AddSource(s)
}

// Register wraps around Manager.Register and can be accessed via configparser.Register.
func Register(name, desc string, defaultValue interface{}) (*Options, error) {
	return Standalone.Register(name, desc, defaultValue)
}

// Load wraps around Manager.Load and can be accessed via configparser.Load.
func Load() {
	Standalone.Load()
}
