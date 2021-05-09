package main

import (
	"github.com/Iridaceae/iridaceae/pkg"
	"github.com/Iridaceae/iridaceae/pkg/deprecatedrunner"
	"github.com/Iridaceae/iridaceae/pkg/sclog/log"
)

// depart all deprecatedrunner run into internal.
func main() {

	defer log.Info("--shutdown--")
	// we will handle all flags here

	_ = pkg.LoadGlobalEnv()
	// TODO: should check if it is running inside docker or a CI pipe
	log.Warn("Make sure that envars are set correctly in docker and CI.")

	if err := pkg.LoadConfig(pkg.ConcertinaClientID, pkg.ConcertinaClientSecrets, pkg.ConcertinaBotToken); err != nil {
		log.Error(err, "couldn't load required envars.")
	}
	// setup metrics here.
	// ....

	// Start bot finally.
	ir := deprecatedrunner.New()
	err := ir.Start()
	if err != nil {
		log.Error(err)
	}
}
