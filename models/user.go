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
	ID       bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	Login    string        `json:"login" bson:"login,omitempty"`
	Password string        `json:"password" bson:"password,omitempty"`
	Email    string        `json:"email" bson:"email,omitempty"`
	Payload  string        `json:"payload,omitempty" bson:"payload,omitempty"`
}

// Implement stringer
func (u User) String() string {
	return fmt.Sprintf("Id: %v\tLogin: %v\tPassword: %v\tEmail: %v\nPayload: %v", u.ID, u.Login, u.Password, u.Email, u.Payload)
}

// NewToken creates a new token using private key.
// pk - the private key.
func (u User) NewToken(pk *rsa.PrivateKey) (string, error) {
	token := jwt.New(jwt.SigningMethodRS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["email"] = u.Email
	claims["login"] = u.Login
	claims["iat"] = time.Now().Unix()

	rawTokenString, err := token.SigningString()
	if err != nil {
		return "", err
	}

	sign, err := token.Method.Sign(rawTokenString, pk)
	if err != nil {
		return "", err
	}

	return rawTokenString + "." + sign, nil
}
