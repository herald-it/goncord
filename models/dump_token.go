package models

import (
	"gopkg.in/mgo.v2/bson"
)

type dumpToken struct {
	UserId bson.ObjectId `json:"user_id" bson:"user_id"`
	Token  string        `json:"token" bson:"token"`
}

// Create new dump token struct.
func NewDumpToken(u *User, token string) dumpToken {
	return dumpToken{
		UserId: u.Id,
		Token:  token,
	}
}
