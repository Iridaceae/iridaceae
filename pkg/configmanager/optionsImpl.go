package configmanager

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// optionsImpl is our variable configuration.
type optionsImpl struct {
	Name         string // given name format iris.option1.option2.
	Description  string
	DefaultValue interface{}
	LoadedValue  interface{}
	Manager      Manager // all manager will implements our default implementation.
	ConfigSource Source
}

func (o *optionsImpl) LoadValue() {
	_default := o.DefaultValue
	man, ok := o.Manager.(*managerImpl)
	// This deals with manager implementation without null source.
	if !ok {
		panic("given option manager doesn't implement iridaceae internal manager")
	}
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
		if _, ok = o.DefaultValue.(int); ok {
			_default = interface{}(toIntVal(_default))
		} else if _, ok = o.DefaultValue.(bool); ok {
			_default = interface{}(toBoolVal(_default))
		}
	}

	o.LoadedValue = _default
}

func (o *optionsImpl) UpdateValue(val interface{}) {
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

func (o optionsImpl) GetValue() interface{} {
	return o.LoadedValue
}

func (o optionsImpl) GetName() string {
	return o.Name
}

func (o *optionsImpl) ToString() string {
	return toStrVal(o.LoadedValue)
}

func (o *optionsImpl) ToInt() int {
	return toIntVal(o.LoadedValue)
}

func (o *optionsImpl) ToBool() bool {
	return toBoolVal(o.LoadedValue)
}

func (o *optionsImpl) ToFloat() float64 {
	return toFloat64Val(o.LoadedValue)
}

func toStrVal(i interface{}) string {
	switch t := i.(type) {
	case string:
		return t
	case int:
		return strconv.FormatInt(int64(t), 10)
	case bool:
		return strconv.FormatBool(t)
	case float64:
		return strconv.FormatFloat(t, 'f', 3, 64)
	case fmt.Stringer:
		return t.String()
	default:
		panic("cannot convert given input to string")
	}
}

func toIntVal(i interface{}) int {
	switch t := i.(type) {
	case string:
		if n, ok := strconv.Atoi(t); ok == nil {
			return n
		}
		return 0
	case int, float64:
		return t.(int)
	case bool:
		if t {
			return 1
		}
		return 0
	default:
		panic("cannot convert given input to int")
	}
}

func toFloat64Val(i interface{}) float64 {
	switch t := i.(type) {
	case string:
		n, ok := strconv.ParseFloat(t, 64)
		if ok == nil {
			return n
		}
		return 0
	case int:
		return float64(t)
	case float64:
		return t
	default:
		panic("cannot convert given input to float64")
	}
}

func toBoolVal(i interface{}) bool {
	switch t := i.(type) {
	case string:
		lower := strings.ToLower(strings.TrimSpace(t))
		if lower == "true" || lower == "yes" || lower == "on" || lower == "enabled" || lower == "1" {
			return true
		}
		return false
	case int, float64:
		return t.(int) > 0
	case bool:
		return t
	default:
		panic("cannot convert given input to bool")
	}
}

func matchOptionsRegex(key string) (bool, error) {
	b, _ := regexp.MatchString(OptionsRegex, key)
	if b {
		return b, nil
	}
	return b, ErrInvalidOptionsMatch
}
