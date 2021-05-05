package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Iridaceae/iridaceae/pkg"

	"github.com/bwmarrin/discordgo"

	"github.com/joho/godotenv"

	"github.com/Iridaceae/iridaceae/internal"
)

const defaultConfigPath = "./defaults.env"

func main() {
	logger := internal.NewLogger(internal.Debug, "concertina")
	defer logger.Info("--shutdown--")

	// parse configparser and secrets parent directory since viper will handle configparser.
	configPath := flag.String("configPath", defaultConfigPath, fmt.Sprintf("LogConfig path for storing default configparser and secrets, default: %s", defaultConfigPath))
	// NOTE: this is when parsing options to get metrics from prom.
	// var opts metricsOptions

	flag.Parse()

	err := godotenv.Load(*configPath)
	if err != nil {
		logger.Warn(fmt.Sprintf("Error loading env file: %s, loading from ENVARS instead.", err.Error()))
	}

	err = pkg.LoadConfig(pkg.ConcertinaClientID, pkg.ConcertinaClientSecrets, pkg.ConcertinaBotToken)
	if err != nil {
		return
	}
	dg, err := discordgo.New(pkg.GetBotToken(pkg.ConcertinaBotToken))
	if err != nil {
		panic(err)
	}
	err = dg.Open()
	if err != nil {
		panic(err)
	}

	logger.Info("Concertina is now running. Press CTRL-C to exit.")
	// Wait for the user to cancel the process.
	defer func() {
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-sc
	}()
}
