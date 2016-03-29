package models

import (
	"fmt"
	"gopkg.in/mgo.v2"
)

type User struct {
	Login    string `json:"login" bson:"login"`
	Password string `json:"password" bson:"password"`
	Email    string `json:"email" bson:"email"`
}

// Implement stringer
func (u User) String() string {
	return fmt.Sprintf("Login: %v\tPassword: %v\tEmail: %v\n", u.Login, u.Password, u.Email)
}

func (u User) Find(c *mgo.Collection) *mgo.Query {
	find_user_query := `$or: [
        {"login": %v},
        {"email": %v}
    ]`

	query := fmt.Sprintf(find_user_query, u.Login, u.Email)
	collect := c.Find(query)

	return collect
}
