// Package configmanager defines some default configs handler including a configparser parser with ability to update value dynamically
package configmanager

// managerImpl holds types for generic managers to generate configs.
type managerImpl struct {
	sources []Source
	Options map[string]*optionsImpl
}

// NewDefaultManager makes a configs manager.
func NewDefaultManager() Manager {
	return &managerImpl{Options: make(map[string]*optionsImpl)}
}

func (c *managerImpl) RegisterSource(source Source) {
	c.sources = append(c.sources, source)
}

func (c *managerImpl) RegisterOption(name, desc string, defaultValue interface{}) (Options, error) {
	if _, err := matchOptionsRegex(name); err != nil {
		return nil, ErrInvalidFormat
	}
	opt := &optionsImpl{
		Name:         name,
		Description:  desc,
		DefaultValue: defaultValue,
		Manager:      c,
	}
	c.Options[name] = opt
	return opt, nil
}

func (c *managerImpl) LoadOptions() {
	for _, v := range c.Options {
		v.LoadValue()
	}
}

func (c *managerImpl) Clear(source bool, option bool) {
	if source {
		c.sources = make([]Source, 0)
	}
	if option {
		c.Options = make(map[string]*optionsImpl)
	}
}
