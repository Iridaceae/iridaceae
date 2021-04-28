// Package dbstore handles all database configuration
package dbstore

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"

	"github.com/globalsign/mgo/bson"
)

func init() {
	err := godotenv.Load("./deployments/defaults.env")
	if err != nil {
		logger.Warn("Error loading env file: %s", err.Error())
	}

	mUser := os.Getenv("MONGO_USER")
	mPass := os.Getenv("MONGO_PASS")
	mIP := parseMongoAddr(os.Getenv("MONGO_ADDR"))

	initMgoSessions(mUser, mPass, mIP)
}

// uri format: this-shard-00-00.asdfg.mongodb.net:27017,this-shard-00-01.asdfg.mongodb.net:27017,this-shard-00-02.asdfg.mongodb.net:27017.
func parseMongoAddr(uri string) []string {
	return strings.Split(uri, ",")
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
