package configmanager

import (
	"os"
	"strings"
)

// EnvSource defines source for environment variables.
type EnvSource struct{}

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

func (e *EnvSource) Name() string {
	return "ENV"
}
