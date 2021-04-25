package main

import (
	"flag"
	"fmt"
	"os"

	bela "github.com/TensRoses/iris/belamcanda"
	"github.com/TensRoses/iris/internal/config"
	"github.com/TensRoses/iris/internal/log"
)

const (
	defaultConfigPath = "./envars"
)

// type metricsOptions struct {
// 	PrometheusMetrics  bool
// 	PrintMetrics       bool
// 	StackDriverMetrics bool
// 	StatsdMetrics      bool
// }

// depart all core run into internal.
func main() {
	logger := log.CreateLogger("tensrose")
	defer logger.Infof("--shutdown %s--", logger.Name)

	// parse config and secrets parent directory since viper will handle config
	cpath := flag.String("cpath", defaultConfigPath, fmt.Sprintf("Config path for storing default config and secrets, default: %s", defaultConfigPath))
	// NOTE: this is when parsing options to get metrics from prom
	// var opts metricsOptions

	flag.Parse()

	// load config and secrets
	cfg, err := config.LoadConfigFile(*cpath)
	if err != nil {
		logger.Fatal(err)
	}

	// NOTE: possible caveats with multiple instance of viper
	// https://stackoverflow.com/a/47185439/8643197
	secret, err := config.LoadSecretsFile(*cpath)
	if err != nil {
		logger.Warnf("%s, loading from ENV instead", err.Error())
		secret.AuthToken = os.Getenv("AUTH_TOKEN")
		secret.ClientID = os.Getenv("CLIENTID")
	}

	// setup metrics here
	// ....

	// Start bela finally
	ir := bela.NewIris(*cfg, *secret, *logger)
	err = ir.Start()
	if err != nil {
		logger.Fatal(err)
	}
}
