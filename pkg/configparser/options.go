package configparser

// Options is our variable configuration.
type Options struct {
	// Name will have format iris.option1.option2
	Name         string
	Description  string
	DefaultValue interface{}
	LoadedValue  interface{}
	Manager      *ConfigManager
	ConfigSource Source
}

// LoadValue will load given values if exists, otherwise use default ones.
func (o *Options) LoadValue() {
	def := o.DefaultValue
	o.ConfigSource = nil

	for i := len(o.Manager.sources) - 1; i >= 0; i-- {
		source := o.Manager.sources[i]
		// v would be value from given source, check envsource.go for examples
		v, _ := source.GetValue(o.Name)

		if v != nil {
			def = v
			o.ConfigSource = source
			break
		}
	}

	if o.DefaultValue != nil {
		if _, ok := o.DefaultValue.(int); ok {
			def = interface{}(toIntVal(def))
		} else if _, ok = o.DefaultValue.(bool); ok {
			def = interface{}(toBoolVal(def))
		}
	}

	o.LoadedValue = def
}

// UpdateValue updates loaded value.
func (o *Options) UpdateValue(val interface{}) {
	switch val.(type) {
	case bool:
		o.LoadedValue = toBoolVal(val)
	case string:
		o.LoadedValue = toStrVal(val)
	case int:
		o.LoadedValue = toIntVal(val)
	case float64:
		o.LoadedValue = toFloat64Val(val)
	}
}

// GetString are a getter string for &Options.LoadedValue.
func (o *Options) GetString() string {
	return toStrVal(o.LoadedValue)
}

// GetInt are a getter int for &Options.LoadedValue.
func (o *Options) GetInt() int {
	return toIntVal(o.LoadedValue)
}

// GetBool are a getter bool for &Options.LoadedValue.
func (o *Options) GetBool() bool {
	return toBoolVal(o.LoadedValue)
}

// GetFloat are a getter float64 for &Options.LoadedValue.
func (o *Options) GetFloat() float64 {
	return toFloat64Val(o.LoadedValue)
}
