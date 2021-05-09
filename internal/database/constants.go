package datastore

import (
	"github.com/globalsign/mgo"
)

var (
	// Session represents a mgo connection.
	Session *mgo.Session
	users   *mgo.Collection
	// fmt given: mongodb://app:password_here@shard:27017,another-shard:27017.
	uriFmt = "mongodb://%s:%s@%s"
)
