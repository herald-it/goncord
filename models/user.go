package models

type User struct {
	Login    string `json:"login" bson:"login"`
	Password string `json:"password" bson:"password"`
	Email    string `json:"email" bson:"email"`
}
