package models

import (
	"gopkg.in/mgo.v2/bson"
)

type DumpToken struct {
	UserId bson.ObjectId `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Token  string        `json:"token" bson:"token"`
}

// Create new dump token struct.
func NewDumpToken(u *User, token string) DumpToken {
	return DumpToken{
		UserId: u.Id,
		Token:  token,
	}
}
