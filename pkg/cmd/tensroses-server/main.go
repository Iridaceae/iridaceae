package main

import (
	"flag"
	"fmt"

	"github.com/joho/godotenv"

	"github.com/TensRoses/iris/pkg/core"

	"github.com/TensRoses/iris/internal/irislog"
)

const (
	defaultConfigPath = "./defaults.env"
)

// depart all core run into internal.
func main() {
	logger := irislog.NewLogger(irislog.Debug, "tensroses-server").Set()
	defer logger.Info("--shutdown--")

	// parse configparser and secrets parent directory since viper will handle configparser.
	cpath := flag.String("cpath", defaultConfigPath, fmt.Sprintf("Config path for storing default configparser and secrets, default: %s", defaultConfigPath))
	// NOTE: this is when parsing options to get metrics from prom.
	// var opts metricsOptions

	flag.Parse()

	err := godotenv.Load(*cpath)
	if err != nil {
		logger.Warn(fmt.Sprintf("Error loading env file: %s, loading from ENVARS instead.", err.Error()))
	}

	// setup metrics here.
	// ....

	// Start bot finally.
	ir := core.New()
	err = ir.Start()
	if err != nil {
		logger.Error(err)
	}
}
