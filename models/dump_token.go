package models

import (
	"gopkg.in/mgo.v2/bson"
)

type DumpToken struct {
	UserId bson.ObjectId `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Token  string        `json:"token" bson:"token"`
}

// NewDumpToken create new dump token struct.
// Associating a user id with a token.
func NewDumpToken(u *User, token string) DumpToken {
	return DumpToken{
		UserId: u.ID,
		Token:  token,
	}
}
