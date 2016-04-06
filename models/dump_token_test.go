package models_test

import (
	"encoding/json"
	"github.com/herald-it/goncord/models"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewDumpTokenModel(t *testing.T) {
	dump_token := &models.DumpToken{}

	Convey("Create new dump token", t, func() {
		So(dump_token, ShouldNotBeNil)
	})
}

func TestJsonDumpTokenModel(t *testing.T) {
	dump_token := models.DumpToken{
		Token: "my_secret_token"}

	const str = `{"token":"my_secret_token"}`
	b, e := json.Marshal(&dump_token)

	Convey("Marshal struct to json", t, func() {
		So(e, ShouldBeNil)
	})

	Convey("Test correct jsonify", t, func() {
		So(string(b), ShouldEqual, str)
	})
}
