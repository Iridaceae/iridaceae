package main

import (
	"os"
	"os/signal"
	"syscall"

	helpers2 "github.com/Iridaceae/iridaceae/internal/helpers"

	"github.com/Iridaceae/iridaceae/pkg/log"

	"github.com/bwmarrin/discordgo"
)

func main() {
	// Remember to setup our logs first.
	log.Mapper().SetAbsent("name", "concertina")
	log.SetGlobalFields([]string{"name"})
	defer log.Info().Msg("--shutdown--")

	// we will handle all flags here

	_ = helpers2.LoadGlobalEnv()
	// TODO: should check if it is running inside docker or a CI pipe
	log.Warn().Msg("Make sure that envars are set correctly in docker and CI.")

	if err := helpers2.LoadConfig(helpers2.ConcertinaClientID, helpers2.ConcertinaClientSecrets, helpers2.ConcertinaBotToken); err != nil {
		log.Error(err).Msg("couldn't load required envars.")
	}
	dg, err := discordgo.New(helpers2.GetBotToken(helpers2.ConcertinaBotToken))
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
