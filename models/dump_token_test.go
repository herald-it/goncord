package models_test

import (
	"encoding/json"
	"github.com/herald-it/goncord/models"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewDumpTokenModel(t *testing.T) {
	Convey("Create new dump token", t, func() {
		Convey("Create new dump token without constructor", func() {
			dump_token := &models.DumpToken{}
			So(dump_token, ShouldNotBeNil)

		})

		Convey("Create new dump token constructor", func() {
			usr := models.User{}
			dump_token := models.NewDumpToken(&usr, "token")
			So(dump_token.Token, ShouldEqual, "token")
		})
	})
}

func TestJsonDumpTokenModel(t *testing.T) {
	Convey("Model to json format", t, func() {
		dump_token := models.DumpToken{
			Token: "my_secret_token"}

		const str = `{"token":"my_secret_token"}`
		b, e := json.Marshal(&dump_token)

		Convey("Marshal struct to json", func() {
			So(e, ShouldBeNil)
		})

		Convey("Test correct jsonify", func() {
			So(string(b), ShouldEqual, str)
		})
	})
}
