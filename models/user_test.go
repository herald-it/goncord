package models_test

import (
	"encoding/json"
	"github.com/herald-it/goncord/keygen"
	"github.com/herald-it/goncord/models"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewUserModel(t *testing.T) {
	Convey("Create new user", t, func() {
		usr := &models.User{
			Login:    "tl",
			Password: "tp",
			Email:    "te"}

		So(usr, ShouldNotBeNil)
	})
}

func TestJsonUserModel(t *testing.T) {
	Convey("Model to json format", t, func() {
		usr := models.User{
			Login:    "log",
			Password: "pwd",
			Email:    "ema"}

		const str = `{"login":"log","password":"pwd","email":"ema"}`
		b, e := json.Marshal(&usr)

		Convey("Marshal struct to json", func() {
			So(e, ShouldBeNil)
		})

		Convey("Test correct jsonify", func() {
			So(string(b), ShouldEqual, str)
		})
	})
}

func TestNewTokenMethod(t *testing.T) {
	Convey("Test new token", t, func() {
		usr := models.User{
			Login:    "log",
			Password: "pwd",
			Email:    "ema"}

		rsa_key, err := keygen.NewKeyPair()

		Convey("New key pair", func() {
			So(err, ShouldBeNil)
			So(rsa_key, ShouldNotBeNil)
		})

		token, err := usr.NewToken(rsa_key.Private)
		Convey("Create new user token", func() {
			So(err, ShouldBeNil)
			So(token, ShouldNotBeNil)
			So(token, ShouldNotEqual, "")
		})
	})
}
