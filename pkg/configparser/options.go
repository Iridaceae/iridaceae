package configparser

// Options is our variable configuration.
type Options struct {
	Name         string // given name format iris.option1.option2.
	Description  string
	DefaultValue interface{}
	LoadedValue  interface{}
	Manager      Manager // all manager will implements our default implementation.
	ConfigSource Source
}

// LoadValue will load given values if exists, otherwise use default ones.
func (o *Options) LoadValue() {
	_default := o.DefaultValue
	man := o.Manager.(*managerImpl)
	o.ConfigSource = nil

	for i := len(man.sources) - 1; i >= 0; i-- {
		source := man.sources[i]
		v, _ := source.GetValue(o.Name)

		if v != nil {
			_default = v
			o.ConfigSource = source
			break
		}
	}

	if o.DefaultValue != nil {
		if _, ok := o.DefaultValue.(int); ok {
			_default = interface{}(toIntVal(_default))
		} else if _, ok = o.DefaultValue.(bool); ok {
			_default = interface{}(toBoolVal(_default))
		}
	}

	o.LoadedValue = _default
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
