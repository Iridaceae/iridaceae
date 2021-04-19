package models

import "github.com/globalsign/mgo/bson"

// UserSchema defined a user info with stats
type User struct {
	Id             bson.ObjectId `bson:"_id,omitempty"`
	discordID      string
	discordTag     string
	guidID         string
	channelID      string
	minutesStudied int
}
