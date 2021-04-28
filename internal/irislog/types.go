package irislog

import (
	"github.com/TensRoses/iris/internal/configparser"
	"github.com/rs/zerolog"
)

var (
	logCfg     LogConfig
	logmanager = configparser.NewManager()
	LogLevel   = logmanager.Register("irislog.level", "logger level for iris", nil)
	Configured = logmanager.Register("irislog.configured", "boolean for configuration", false)
)

type IrisLogger struct {
	Level    int
	Name     string
	Version  string
	Revision string
	StdLog   zerolog.Logger
	ErrLog   zerolog.Logger
}

type LogConfig struct {
	manager   *configparser.Manager
	logger    *IrisLogger
	stdFields []interface{}
}

func setup(logLevel int, stFields []interface{}) {
	logCfg.manager = logmanager

	Configured.UpdateValue(true)
	LogLevel.UpdateValue(logLevel)
	logCfg.stdFields = append(logCfg.stdFields, stFields...)

	logCfg.manager.Load()
}
