// Package database handles all database configuration
package database

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Iridaceae/iridaceae/pkg/helpers"

	"github.com/Iridaceae/iridaceae/pkg/log"

	"github.com/joho/godotenv"

	"github.com/globalsign/mgo/bson"
)

func init() {
	err := godotenv.Load(strings.Join([]string{helpers.GetRootDir(), "defaults.env"}, "/"))
	if err != nil {
		log.Error(err).Msg("Error loading env file")
	}

	mUser := os.Getenv("IRIS_MONGO_USER")
	mPass := os.Getenv("IRIS_MONGO_PASS")
	mDBName := os.Getenv("IRIS_MONGO_DBNAME")
	mIP := os.Getenv("IRIS_MONGO_ADDR")
	mAddr := fmt.Sprintf(uriFmt, mUser, mPass, mIP)

	initMgoSessions(mDBName, mAddr)
}

// NewUser returns a hex representation of the inputs ObjectID and insert errors into new database.
func NewUser(did, dit, guild, minutes string) (string, error) {
	m, _ := strconv.Atoi(minutes)
	oid := bson.NewObjectId()
	newEntry := User{
		ID:             oid,
		DiscordID:      did,
		DiscordTag:     dit,
		GUILDID:        guild,
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
func UpdateUser(did, guildid, channelid string, minutes int) error {
	return update(did, guildid, channelid, minutes)
}

// FetchNumHours returns total number of hours of given users.
func FetchNumHours(did string) string {
	u, err := fetch(did)
	if err != nil {
		log.Error(err).Msgf("error while fetching users %s", did)
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
