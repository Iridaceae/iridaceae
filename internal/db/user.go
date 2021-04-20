package db

import "github.com/globalsign/mgo/bson"

// UserSchema defined a user info with stats
type User struct {
	Id             bson.ObjectId `bson:"_id,omitempty"`
	DiscordID      string        `bson:"discordid"`
	DiscordTag     string        `bson:"discordtag"`
	GuidID         string        `bson:"guidid"`
	ChannelID      string        `bson:"channelid"`
	MinutesStudied int           `bson:"minutesstudied"`
}
