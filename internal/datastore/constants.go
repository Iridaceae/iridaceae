package datastore

import (
	"os/exec"

	"github.com/globalsign/mgo"

	"github.com/TensRoses/iris/internal/irislog"
)

const name string = "dbstore_service"

var (
	rev     = getRevision()
	logger  = irislog.NewLogger(irislog.Debug, name, "revision", rev)
	Session *mgo.Session
	users   *mgo.Collection
)

func getRevision() string {
	// check for errors instead of printing to os.Stdout
	stdout, _ := exec.Command("git", "rev-parse", "--short", "HEAD").Output()
	return string(stdout)
}
