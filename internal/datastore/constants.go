package datastore

import (
	"os/exec"

	"github.com/Iridaceae/iridaceae/internal"

	"github.com/globalsign/mgo"
)

const name string = "dbstore_service"

var (
	rev      = getRevision()
	dbLogger = internal.NewLogger(internal.Debug, name, "revision", rev).Set()
	// Session represents a mgo connection.
	Session *mgo.Session
	users   *mgo.Collection
	// fmt given: mongodb://app:password_here@shard:27017,another-shard:27017.
	uriFmt = "mongodb://%s:%s@%s"
)

func getRevision() string {
	// check for errors instead of printing to os.Stdout
	stdout, _ := exec.Command("git", "rev-parse", "--short", "HEAD").Output()
	return string(stdout)
}
