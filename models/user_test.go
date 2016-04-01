package models_test

import (
	"github.com/herald-it/goncord/models"
	"testing"
)

func TestNewUserModel(t *testing.T) {
	usr := models.User{
		Login:    "tl",
		Password: "tp",
		Email:    "te"}

	if false {
		usr = usr
	}
}
