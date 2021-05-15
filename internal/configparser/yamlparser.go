package configparser

import (
	"io"

	"gopkg.in/yaml.v3"
)

type YamlParser struct{}

func (y *YamlParser) Marshal(w io.Writer, s *Settings) error {
	b, err := yaml.Marshal(s)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}

func (y *YamlParser) Unmarshal(r io.Reader) (*Settings, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	c := new(Settings)
	err = yaml.Unmarshal(b, c)
	return c, err
}
