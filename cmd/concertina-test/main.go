package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Iridaceae/iridaceae/pkg/log"

	"github.com/Iridaceae/iridaceae/pkg"

	"github.com/bwmarrin/discordgo"
)

func main() {
	// Remember to setup our logs first.
	log.Mapper().SetAbsent("name", "concertina")
	defer log.Info().Msg("--shutdown--")

	// we will handle all flags here

	_ = pkg.LoadGlobalEnv()
	// TODO: should check if it is running inside docker or a CI pipe
	log.Warn().Msg("Make sure that envars are set correctly in docker and CI.")

	if err := pkg.LoadConfig(pkg.ConcertinaClientID, pkg.ConcertinaClientSecrets, pkg.ConcertinaBotToken); err != nil {
		log.Error(err).Msg("couldn't load required envars.")
	}
	dg, err := discordgo.New(pkg.GetBotToken(pkg.ConcertinaBotToken))
	if err != nil {
		panic(err)
	}
	_ = dg.Open()

	log.Info().Msg("Running. Press CTRL-C to exit.")
	// Wait for the user to cancel the process.
	defer func() {
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-sc
	}()
}
