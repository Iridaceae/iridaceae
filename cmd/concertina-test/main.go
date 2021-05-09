package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Iridaceae/iridaceae/pkg/sclog"

	"github.com/Iridaceae/iridaceae/pkg/sclog/log"

	"github.com/Iridaceae/iridaceae/pkg"

	"github.com/bwmarrin/discordgo"
)

func main() {
	defer log.Info("--shutdown--")

	// we will handle all flags here

	_ = pkg.LoadGlobalEnv()
	// TODO: should check if it is running inside docker or a CI pipe
	log.Warn("Make sure that envars are set correctly in docker and CI.")

	if err := pkg.LoadConfig(pkg.ConcertinaClientID, pkg.ConcertinaClientSecrets, pkg.ConcertinaBotToken); err != nil {
		log.Error(err, "couldn't load required envars.")
	}
	dg, err := discordgo.New(pkg.GetBotToken(pkg.ConcertinaBotToken))
	if err != nil {
		panic(err)
	}
	_ = dg.Open()

	sclog.Mapper().Set("name", "concertina")
	sclog.AddGlobalFields("name")
	log.Info("Running. Press CTRL-C to exit.")
	// Wait for the user to cancel the process.
	defer func() {
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-sc
	}()
}
