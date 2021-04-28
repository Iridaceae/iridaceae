package main

import (
	"flag"
	"fmt"

	"github.com/joho/godotenv"

	"github.com/TensRoses/iris/pkg/belamcanda"

	"github.com/TensRoses/iris/internal/irislog"
)

const (
	defaultConfigPath = "./deployments/defaults.env"
)

// type metricsOptions struct {
// 	PrometheusMetrics  bool
// 	PrintMetrics       bool
// 	StackDriverMetrics bool
// 	StatsdMetrics      bool
// }

// depart all core run into internal.
func main() {
	logger := irislog.NewLogger(irislog.Debug, "tensroses-server")
	defer logger.Info("--shutdown--")

	// parse configparser and secrets parent directory since viper will handle configparser
	cpath := flag.String("cpath", defaultConfigPath, fmt.Sprintf("Config path for storing default configparser and secrets, default: %s", defaultConfigPath))
	// NOTE: this is when parsing options to get metrics from prom
	// var opts metricsOptions

	flag.Parse()

	err := godotenv.Load(*cpath)
	if err != nil {
		logger.Warn("Error loading env file: %s", err.Error())
	}

	// setup metrics here
	// ....

	// Start bela finally
	ir := belamcanda.NewIris()
	err = ir.Start()
	if err != nil {
		logger.Error(err)
	}
}
