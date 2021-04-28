package dbstore

import "github.com/globalsign/mgo/bson"

// User defined a user info with stats.
type User struct {
	Id             bson.ObjectId `bson:"_id,omitempty"`
	DiscordId      string        `bson:"discordid"`
	DiscordTag     string        `bson:"discordtag"`
	GuidId         string        `bson:"guidid"`
	ChannelId      string        `bson:"channelid"`
	MinutesStudied int           `bson:"minutesstudied"`
}
