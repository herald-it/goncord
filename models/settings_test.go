package models_test

import (
	"testing"

	"github.com/herald-it/goncord/models"
	. "github.com/smartystreets/goconvey/convey"
)

func TestLoadSettings(t *testing.T) {
	Convey("Load settings from test", t, func() {
		err := models.LoadSettings()

		So(err, ShouldBeNil)
		So(models.Set.Database, ShouldNotBeNil)
		So(models.Set.Database.DbName, ShouldNotBeNil)
		So(models.Set.Database.Host, ShouldNotBeNil)
		So(models.Set.Database.TokenTable, ShouldNotBeNil)
		So(models.Set.Database.UserTable, ShouldNotBeNil)

		Convey("Test valid parse setting file", func() {
			Convey("Database section", func() {
				So(models.Set.Database.Host, ShouldEqual, "localhost")
				So(models.Set.Database.DbName, ShouldEqual, "testdb")
				So(models.Set.Database.TokenTable, ShouldEqual, "tokentable")
				So(models.Set.Database.UserTable, ShouldEqual, "usertable")
			})
			Convey("Ssl section", func() {
				So(models.Set.Ssl.Key, ShouldEqual, "./pass_key")
				So(models.Set.Ssl.Certificate, ShouldEqual, "./pass_certificate")
			})
			Convey("Router section", func() {
				So(models.Set.Router.Login, ShouldEqual, "log")
				So(models.Set.Router.Register, ShouldEqual, "reg")
				So(models.Set.Router.Validate, ShouldEqual, "valid")
			})

			So(models.Set.Domain, ShouldEqual, "my.domain.com")
			So(models.Set.IP, ShouldEqual, "0.0.0.0:8000")
		})
	})
}
