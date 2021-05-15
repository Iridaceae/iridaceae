package configmanager

import (
	"io"

	"github.com/Iridaceae/iridaceae/internal/components"

	"gopkg.in/yaml.v3"
)

type YamlParser struct{}

func (y *YamlParser) Marshal(w io.Writer, c *Config) error {
	b, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}

func (y *YamlParser) Unmarshal(r io.Reader) (*Config, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	c := new(Config)
	err = yaml.Unmarshal(b, c)
	return c, err
}

func GetIrisConfig() *Config {
	return &Config{
		Version:    1,
		Discord:    &Discord{GlobalRateLimit: &GlobalRateLimit{Burst: 2, Restoration: 10}},
		Permission: &Permission{User: components.DefaultUserRules, Admin: components.DefaultAdminRules},
		Logging:    &Logging{Enabled: true, Level: 1},
		Metrics:    &Metrics{Enabled: true, Address: ":9000"},
	}
}
