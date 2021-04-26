package configparser

import (
	"os"
	"strings"
)

type EnvSource struct{}

// GetValue will get env vars with following format for config parsing: iris.option1.option2
func (e *EnvSource) GetValue(key string) interface{} {
	envKey := strings.ToUpper(key)
	envKey = strings.ReplaceAll(envKey, ".", "_")
	v := os.Getenv(envKey)
	if v == "" {
		return nil
	}
	return v
}

func (e *EnvSource) Name() string  {
	return "ENV"
}