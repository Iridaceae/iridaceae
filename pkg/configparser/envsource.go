package configparser

import (
	"os"
	"strings"
)

// EnvSource defines source for environment variables.
type EnvSource struct{}

// GetValue will get env vars with following format for config parsing: iris.option1.option2.
// and will returns keys as follow IRIS_OPTION1_OPTION2.
func (e *EnvSource) GetValue(key string) (interface{}, error) {
	b, _ := matchOptionsRegex(key)

	if b {
		envKey := strings.ToUpper(key)
		// NOTE: add options to check for correct options parsing
		envKey = strings.ReplaceAll(envKey, ".", "_")
		v := os.Getenv(envKey)
		if v == "" {
			return nil, ErrEmptyValue
		}
		return v, nil
	}
	return nil, ErrInvalidFormat
}

// Name will return the name of given source.
func (e *EnvSource) Name() string {
	return "ENV"
}
