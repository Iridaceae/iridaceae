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
	err := godotenv.Load("./defaults.env")
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
		ID:             oid,
		DiscordID:      did,
		DiscordTag:     dit,
		GUIDID:         guid,
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

// FetchNumHours returns total number of hours of given users.
func FetchNumHours(did string) string {
	u, err := fetch(did)
	if err != nil {
		dbLogger.Warn("err", fmt.Sprintf("error while fetching users %s: %s", did, err.Error()))
	}
	return toHumanTime(u.MinutesStudied)
}

// since time is captured in minutes, it will omit the following format 1d2h4m.
func toHumanTime(time int) string {
	day := time / 60 / 24
	hours := time / 60 % 24
	minutes := time % 60
	if day > 0 {
		return fmt.Sprintf("%dd %dh %dm", day, hours, minutes)
	}
	return fmt.Sprintf("%dh %dm", hours, minutes)
}
