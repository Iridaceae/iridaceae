// mongo.go contains DAO operations that can later be called with your repository layer.
// DAO lives within mongo.go for now, as per implementation only requires it to be as mongo. SQL will be used for logging and metrics
package db

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"

	"github.com/aarnphm/iris/internal/log"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

var (
	Session *mgo.Session
	users   *mgo.Collection
	logger  *log.Logging
)

func initMgoSessions(user, pass, ip, port, dbname string) {
	logger.Info("attempting to connect at %s", ip)

	// building URI
	URIfmt := "mongodb://%s:%s@%s:%s/%s"
	// URI from mongo mongodb+srv://admin:<password>@main.hzdkk.mongodb.net/main?retryWrites=true&w=majority
	MongoURI := fmt.Sprintf(URIfmt, user, pass, ip, port, dbname)

	// https://stackoverflow.com/a/42522753/8643197
	dialInfo, err := mgo.ParseURL(MongoURI)
	if err != nil {
		logger.Fatal(errors.New(fmt.Sprintf("errors parsing uri: %s", err.Error())))
	}

	tlsConfig := &tls.Config{}
	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		return conn, err
	}

	Session, err = mgo.DialWithInfo(dialInfo)
	if err != nil {
		logger.Fatal(errors.New(fmt.Sprintf("error while establishing connection with mongo: %s", err.Error())))
	}

	users = Session.DB("main").C("users")
}

func insert(user User) error {
	return users.Insert(user)
}

func update(discordID string, mins int) error {

	u, err := fetch(discordID)
	if err != nil {
		logger.Fatal(errors.New(fmt.Sprintf("error while finding discordID: %s", err.Error())))
	}

	newMin := u.MinutesStudied + mins

	minStudy := bson.M{"minutesstudied": newMin}
	err = users.Update(minStudy, bson.M{"$set": u})
	return err
}

func fetch(discordID string) (User, error) {
	u := User{}

	uid := bson.M{"discordid": discordID}
	err := users.Find(uid).One(&u)
	return u, err
}
