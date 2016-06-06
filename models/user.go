package models

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2/bson"
)

// User model.
type User struct {
	Id                  bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty" schema:"-"`
	Login               string        `json:"login" bson:"login,omitempty" scheme:"login"`
	Password            string        `json:"password" bson:"password,omitempty" scheme:"password"`
	Email               string        `json:"email" bson:"email,omitempty" scheme:"email"`
	Csrfmiddlewaretoken string        `scheme:"csrfmiddlewaretoken"`
}

// Implement stringer
func (u User) String() string {
	return fmt.Sprintf("Id: %v\tLogin: %v\tPassword: %v\tEmail: %v\n", u.Id, u.Login, u.Password, u.Email)
}

// NewToken creates a new token using private key.
// pk - the private key.
func (u User) NewToken(pk *rsa.PrivateKey) (string, error) {
	token := jwt.New(jwt.SigningMethodRS256)
	token.Claims["email"] = u.Email
	token.Claims["login"] = u.Login
	token.Claims["iat"] = time.Now().Unix()

	rawTokenStrng, err := token.SigningString()
	if err != nil {
		return "", err
	}

	sign, err := token.Method.Sign(rawTokenStrng, pk)
	if err != nil {
		return "", err
	}

	return rawTokenStrng + "." + sign, nil
}
