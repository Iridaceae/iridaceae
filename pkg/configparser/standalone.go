package configparser

// Standalone is a singleton config manager that acts as a general manager for iris.
var Standalone = NewConfigManager()

// TestParser acts as a test config manager that can be used globally.
var TestParser = NewConfigManager()

func AddSource(s Source) {
	Standalone.AddSource(s)
}

func Register(name, desc string, defaultValue interface{}) (*Options, error) {
	return Standalone.Register(name, desc, defaultValue)
}

func Load() {
	Standalone.Load()
}
