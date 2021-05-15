// Package configparser defines some default configs handler including a configparser parser with ability to update value dynamically
package configparser

// managerImpl holds types for generic managers to generate configs.
type managerImpl struct {
	sources []Source
	Options map[string]*Options
}

// NewConfigManager makes a configs manager.
func NewConfigManager() Manager {
	return &managerImpl{Options: make(map[string]*Options)}
}

func (c *managerImpl) AddSource(source Source) {
	c.sources = append(c.sources, source)
}

func (c *managerImpl) Register(name, desc string, defaultValue interface{}) (*Options, error) {
	if _, err := matchOptionsRegex(name); err != nil {
		return nil, ErrInvalidFormat
	}
	opt := &Options{
		Name:         name,
		Description:  desc,
		DefaultValue: defaultValue,
		Manager:      c,
	}
	c.Options[name] = opt
	return opt, nil
}

func (c *managerImpl) Load() {
	for _, v := range c.Options {
		v.LoadValue()
	}
}

func (c *managerImpl) Reset() {
	c.sources = make([]Source, 0)
	c.Options = make(map[string]*Options)
}
