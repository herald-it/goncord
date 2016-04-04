package models_test

import (
	"encoding/json"
	"github.com/herald-it/goncord/keygen"
	"github.com/herald-it/goncord/models"
	"testing"
)

func TestNewUserModel(t *testing.T) {
	usr := &models.User{
		Login:    "tl",
		Password: "tp",
		Email:    "te"}

	if usr == nil {
		t.Fatal("Nil pointer after create new user.")
	}
}

func TestJsonUserModel(t *testing.T) {
	usr := models.User{
		Login:    "log",
		Password: "pwd",
		Email:    "ema"}

	const str = `{"login":"log","password":"pwd","email":"ema"}`
	b, e := json.Marshal(&usr)

	if e != nil {
		t.Fatalf("Error: %v", e.Error())
	}

	if string(b) != str {
		t.Fatalf("%v not equal %v", string(b), str)
	}
}

func TestNewTokenMethod(t *testing.T) {
	usr := models.User{
		Login:    "log",
		Password: "pwd",
		Email:    "ema"}

	rsa_key, err := keygen.NewKeyPair()
	if err != nil {
		t.Fatalf("%v", err.Error())
	}

	token, err := usr.NewToken(rsa_key.Private)
	if err != nil {
		t.Fatalf("%v", err.Error())
	}

	if token == "" {
		t.Fatalf("Empty token: %v", token)
	}
}
