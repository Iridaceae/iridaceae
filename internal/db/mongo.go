// mongo.go contains DAO operations that can later be called with your repository layer.
// DAO lives within mongo.go for now, as per implementation only requires it to be as mongo. SQL will be used for logging and metrics
package db

import (
	"crypto/tls"
	"fmt"
	"net"

	"github.com/globalsign/mgo"

	"github.com/globalsign/mgo/bson"
)

var (
	Session *mgo.Session
	users   *mgo.Collection
)

func initMgoSessions(user, pass, ip string) {
	// building URI
	URIfmt := "mongodb://%s:%s@%s:27017"

	logger.Infof("attempting to connect at "+URIfmt, "<user>", "<pass>", ip)

	mongoURI := fmt.Sprintf(URIfmt, user, pass, ip)

	// https://stackoverflow.com/a/42522753/8643197
	dialInfo, err := mgo.ParseURL(mongoURI)
	if err != nil {
		logger.Fatal(fmt.Errorf("errors parsing URI: %w", err))
	}

	tlsConfig := &tls.Config{}
	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, er := tls.Dial("tcp", addr.String(), tlsConfig)
		return conn, er
	}

	Session, err = mgo.DialWithInfo(dialInfo)
	if err != nil {
		logger.Fatal(fmt.Errorf("error while establishing connection with mongo: %w", err))
	}

	users = Session.DB("main").C("users")
}

func insert(user User) error {
	return users.Insert(user)
}

func fetch(discordID string) (User, error) {
	var u User

	err := users.Find(bson.M{"discordid": discordID}).One(&u)
	return u, err
}

func update(discordID string, mins int) error {
	u, _ := fetch(discordID)

	logger.Infof("users id: %s", discordID)
	logger.Infof("mins to be added: %d", mins)
	var newMin = u.MinutesStudied + mins
	logger.Infof("users current score: %d", u.MinutesStudied)
	logger.Infof("new mins: %d", newMin)

	err := users.Update(bson.M{"minutesstudied": newMin}, bson.M{"$set": &u})
	return err
}
