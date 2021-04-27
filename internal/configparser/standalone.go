package configparser

// Standalone is a singleton config manager that acts as a general manager for iris
var Standalone = NewManager()

// AddSource wraps around Manager.AddSource
func AddSource(s Source){
	Standalone.AddSource(s)
}

// Register wraps around Manager.Register
func Register(name, desc string, defaultValue interface{}) *Options {
	return Standalone.Register(name, desc, defaultValue)
}

// Load wraps around Manager.Load
func Load() {
	Standalone.Load()
}
