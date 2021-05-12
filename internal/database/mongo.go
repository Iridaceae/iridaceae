package database

import (
	"crypto/tls"
	"net"
	"time"

	"github.com/Iridaceae/iridaceae/pkg/log"

	"github.com/globalsign/mgo"

	"github.com/globalsign/mgo/bson"
)

var (
	// Session represents a mgo connection.
	Session *mgo.Session
	users   *mgo.Collection
	// fmt given: mongodb://app:password_here@shard:27017,another-shard:27017.
	uriFmt = "mongodb://%s:%s@%s"
)

// User defined a user info with stats.
type User struct {
	ID             bson.ObjectId `bson:"_id,omitempty"`
	GUILDID        string        `bson:"guildid"`
	DiscordID      string        `bson:"discordid"`
	DiscordTag     string        `bson:"discordtag"`
	ChannelID      string        `bson:"channelid"`
	MinutesStudied int           `bson:"minutesstudied"`
}

func initMgoSessions(dbname, addr string) {
	var err error

	// https://stackoverflow.com/a/42522753/8643197.
	// here we pass addr as replica sets.
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS13,
	}
	dialInfo, _ := mgo.ParseURL(addr)
	dialInfo.Timeout = 5 * time.Second
	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, er := tls.Dial("tcp", addr.String(), tlsConfig)
		return conn, er
	}

	Session, err = mgo.DialWithInfo(dialInfo)
	if err != nil {
		log.Error(err).Msg("error while establishing connection with mongo")
	}

	users = Session.DB(dbname).C("users")
}

func insert(user User) error {
	return users.Insert(user)
}

func fetch(discordID string) (User, error) {
	var u User

	log.Debug().Msgf("fetching %s from db", discordID)
	err := users.Find(bson.M{"discordid": discordID}).One(&u)
	return u, err
}

func update(discordID, guildID, channelID string, mins int) error {
	u, _ := fetch(discordID)

	newMin := u.MinutesStudied + mins

	err := users.Update(bson.M{"discordid": u.DiscordID}, bson.M{"$set": bson.M{"guildid": guildID, "channelid": channelID, "minutesstudied": newMin}})
	return err
}
