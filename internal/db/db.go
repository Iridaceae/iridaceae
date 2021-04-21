package db

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"github.com/globalsign/mgo/bson"

	"github.com/TensRoses/iris/internal/log"
)

var logger *log.Logging = log.CreateLogger("db")

func init() {
	err := godotenv.Load()
	if err != nil {
		logger.Warnf("Error loading env file: %s", err.Error())
	}

	mUser := os.Getenv("MONGO_USER")
	mPass := os.Getenv("MONGO_PASS")
	mIP := os.Getenv("MONGO_ADDR")

	initMgoSessions(mUser, mPass, mIP)
}

// NewUser returns a hex representation of the inputs ObjectID and insert errors into new database.
func NewUser(did, dit, guid, mins string) (string, error) {
	m, _ := strconv.Atoi(mins)
	oid := bson.NewObjectId()
	newEntry := User{
		Id:             oid,
		DiscordID:      did,
		DiscordTag:     dit,
		GuidID:         guid,
		MinutesStudied: m,
	}
	insertErr := insert(newEntry)
	return oid.Hex(), insertErr
}

// UpdateUser updates minutes studied to current users via discordID.
func UpdateUser(did string, mins int) error {
	return update(did, mins)
}

// FetchUser returns a singleton of users given discordID.
func FetchUser(did string) error {
	_, err := fetch(did)
	return err
}
