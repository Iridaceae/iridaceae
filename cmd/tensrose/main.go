package main

import (
	"flag"
	"fmt"

	"github.com/aarnphm/iris/internal/bot"
	"github.com/aarnphm/iris/internal/configs"
	"github.com/aarnphm/iris/internal/log"
)

const (
	defaultConfigPath = "./internal/configs"
)

type metricsOptions struct {
	PrometheusMetrics  bool
	PrintMetrics       bool
	StackdriverMetrics bool
	StatsdMetrics      bool
}

// depart all core run into internal
func main() {
	logger := log.CreateLogger("tensrose")
	defer logger.Info("--SHUTDOWN--")

	// parse configs and secrets parent directory since viper will handle configs
	cpath := flag.String("cpath", defaultConfigPath, fmt.Sprintf("Config path for storing default configs and secrets, default: %s", defaultConfigPath))
	// NOTE: this is when parsing options to get metrics from prom
	// var opts metricsOptions

	flag.Parse()

	// load configs and secrets
	cfg, err := configs.LoadConfigFile(*cpath)
	if err != nil {
		logger.Fatal(err)
	}
	secret, err := configs.LoadSecretsFile(*cpath)
	if err != nil {
		logger.Fatal(err)
	}

	// setup metrics here
	// ....

	//Start Iris finally
	iris := bot.NewIris(*cfg, *secret, *logger)
	err = iris.Start()
	if err != nil {
		logger.Fatal(err)
	}
}
