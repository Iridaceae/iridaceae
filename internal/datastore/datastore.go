// Package datastore handles all database configuration
package datastore

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"github.com/globalsign/mgo/bson"
)

func init() {
	err := godotenv.Load("./deployments/defaults.env")
	if err != nil {
		dbLogger.Warn("Error loading env file: %s", err.Error())
	}

	mUser := os.Getenv("IRIS_MONGO_USER")
	mPass := os.Getenv("IRIS_MONGO_PASS")
	mDBName := os.Getenv("IRIS_MONGO_DBNAME")
	mIP := parseMongoAddr(mUser, mPass, os.Getenv("IRIS_MONGO_ADDR"))

	initMgoSessions(mDBName, mIP)
}

// uri format: this-shard-00-00.asdfg.mongodb.net:27017,this-shard-00-01.asdfg.mongodb.net:27017,this-shard-00-02.asdfg.mongodb.net:27017.
func parseMongoAddr(user, pass, uri string) string {
	return fmt.Sprintf(uriFmt, user, pass, uri)
}

// NewUser returns a hex representation of the inputs ObjectID and insert errors into new database.
func NewUser(did, dit, guid, minutes string) (string, error) {
	m, _ := strconv.Atoi(minutes)
	oid := bson.NewObjectId()
	newEntry := User{
		Id:             oid,
		DiscordId:      did,
		DiscordTag:     dit,
		GuidId:         guid,
		MinutesStudied: m,
	}
	insertErr := insert(newEntry)
	return oid.Hex(), insertErr
}

// FetchUser returns a singleton of users given discordID.
func FetchUser(did string) error {
	_, err := fetch(did)
	return err
}

// UpdateUser updates minutes studied to current users via discordID.
func UpdateUser(did string, minutes int) error {
	return update(did, minutes)
}
