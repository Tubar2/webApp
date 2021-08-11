package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Tweet struct {
	MongoID   primitive.ObjectID `json:"_id"      bson:"_id,omitempty"`
	Author_ID string             `json:"-"        bson:"a_id,omitempty"`
	Tweet     string             `json:"tweet"    bson:"tweet,omitempty"`
}

func NewTweet(a_Id, tweet string) *Tweet {
	return &Tweet{
		MongoID:   primitive.NewObjectID(),
		Author_ID: a_Id,
		Tweet:     tweet,
	}
}

