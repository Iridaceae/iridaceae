package datastore

import (
	"github.com/globalsign/mgo"

	"github.com/Iridaceae/iridaceae/pkg/stlog"
)

const name string = "datastore_service"

var (
	dbLogger = stlog.NewLogger(stlog.Debug, name).Set()
	// Session represents a mgo connection.
	Session *mgo.Session
	users   *mgo.Collection
	// fmt given: mongodb://app:password_here@shard:27017,another-shard:27017.
	uriFmt = "mongodb://%s:%s@%s"
)
