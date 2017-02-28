package models_test

import (
	"testing"

	"github.com/herald-it/goncord/models"
	. "github.com/smartystreets/goconvey/convey"
)

func TestLoadSettings(t *testing.T) {
	Convey("Load settings from test", t, func() {
		err := models.LoadSettings("./settings_test.yml")

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
				So(models.Set.Router.Login.Path, ShouldEqual, "/login")
				So(models.Set.Router.Register.Path, ShouldEqual, "/register")
				So(models.Set.Router.Validate.Path, ShouldEqual, "/validate")
				So(models.Set.Router.Logout.Path, ShouldEqual, "/logout")
				So(models.Set.Router.Update.Path, ShouldEqual, "/update")
				So(models.Set.Router.ResetPassword.Path, ShouldEqual, "/reset")
			})

			So(models.Set.Domain, ShouldEqual, "my.domain.com")
			So(models.Set.IP, ShouldEqual, "0.0.0.0:8000")
			So(models.Set.TelegramToken, ShouldEqual, "test_token")
		})
	})
}
