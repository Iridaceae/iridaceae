package configparser

type YamlSource struct{}

func (y *YamlSource) GetValue(key string) (interface{}, error) {
	return nil, nil
}

func (y *YamlSource) Name() string {
	return "YAML"
}
