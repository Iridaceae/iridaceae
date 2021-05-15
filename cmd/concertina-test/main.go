package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Iridaceae/iridaceae/internal/helpers"

	"github.com/Iridaceae/iridaceae/pkg/log"

	"github.com/bwmarrin/discordgo"
)

func main() {
	// Remember to setup our logs first.
	log.Mapper().SetAbsent("name", "concertina")
	log.SetGlobalFields([]string{"name"})
	defer log.Info().Msg("--shutdown--")

	// we will handle all flags here

	_ = helpers.LoadGlobalEnv()
	// TODO: should check if it is running inside docker or a CI pipe
	log.Warn().Msg("Make sure that envars are set correctly in docker and CI.")

	if err := helpers.LoadConfig(helpers.ConcertinaClientID, helpers.ConcertinaClientSecrets, helpers.ConcertinaBotToken); err != nil {
		log.Error(err).Msg("couldn't load required envars.")
	}
	dg, err := discordgo.New(helpers.GetBotToken(helpers.ConcertinaBotToken))
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
