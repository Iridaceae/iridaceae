package datastore

import "github.com/globalsign/mgo/bson"

// User defined a user info with stats.
type User struct {
	ID             bson.ObjectId `bson:"_id,omitempty"`
	GUIDID         string        `bson:"guidid"`
	DiscordID      string        `bson:"discordid"`
	DiscordTag     string        `bson:"discordtag"`
	ChannelID      string        `bson:"channelid"`
	MinutesStudied int           `bson:"minutesstudied"`
}
