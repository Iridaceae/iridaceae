package configparser

import (
	"io"

	helpers2 "github.com/Iridaceae/iridaceae/internal/helpers"

	"gopkg.in/yaml.v3"
)

type YamlParser struct{}

func (y *YamlParser) Marshal(w io.Writer, c *helpers2.Settings) error {
	b, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}

func (y *YamlParser) Unmarshal(r io.Reader) (*helpers2.Settings, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	c := new(helpers2.Settings)
	err = yaml.Unmarshal(b, c)
	return c, err
}
