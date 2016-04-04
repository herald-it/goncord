package models

import (
	"crypto/rsa"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type User struct {
	Id       bson.ObjectId `json:"_id" bson:"_id,omitempty" schema:"-"`
	Login    string        `json:"login" bson:"login"`
	Password string        `json:"password" bson:"password"`
	Email    string        `json:"email" bson:"email"`
}

// Implement stringer
func (u User) String() string {
	return fmt.Sprintf("Id: %v\tLogin: %v\tPassword: %v\tEmail: %v\n", u.Id, u.Login, u.Password, u.Email)
}

func (u User) NewToken(pk *rsa.PrivateKey) (string, error) {
	token := jwt.New(jwt.SigningMethodRS256)
	token.Claims["email"] = u.Email
	token.Claims["login"] = u.Login
	token.Claims["iat"] = time.Now().Unix()

	raw_token_strng, err := token.SigningString()
	if err != nil {
		return "", err
	}

	sign, err := token.Method.Sign(raw_token_strng, pk)
	if err != nil {
		return "", err
	}

	return raw_token_strng + "." + sign, nil
}
