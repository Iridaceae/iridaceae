package configmanager

// Standalone is a singleton configparser manager that acts as a general manager for iris.
var Standalone = NewDefaultManager().(*managerImpl)

func RegisterSource(s Source) {
	Standalone.RegisterSource(s)
}

func RegisterOption(name, desc string, defaultValue interface{}) (Options, error) {
	return Standalone.RegisterOption(name, desc, defaultValue)
}

func LoadOptions() {
	Standalone.LoadOptions()
}

func Clear() {
	Standalone.Clear(true, true)
}
