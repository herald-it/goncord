package models

import "fmt"

type User struct {
	Login    string `json:"login" bson:"login"`
	Password string `json:"password" bson:"password"`
	Email    string `json:"email" bson:"email"`
}

func (u User) String() string {
	return fmt.Sprintf("Login: %v\tPassword: %v\tEmail: %v\n", u.Login, u.Password, u.Email)
}
