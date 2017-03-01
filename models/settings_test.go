package models_test

import (
	"testing"

	"github.com/herald-it/goncord/models"
	. "github.com/smartystreets/goconvey/convey"
)

func TestLoadSettings(t *testing.T) {
	models.Set = nil

	Convey("Load settings from test", t, func() {
		So(models.Set, ShouldBeNil)
		err := models.LoadSettings("./settings_test.yml")
		So(err, ShouldBeNil)

		Convey("Test valid parse setting file", func() {
			So(models.Set.Domain, ShouldEqual, "my.domain.com")
			So(models.Set.IP, ShouldEqual, "0.0.0.0:8000")
			So(models.Set.Timber.Host, ShouldEqual, "timber_host")
			So(models.Set.Timber.Token, ShouldEqual, "timber_token")
		})
	})
}

func TestRouterSetting(t *testing.T) {
	models.Set = nil

	Convey("Test valid parse setting file", t, func() {
		So(models.Set, ShouldBeNil)
		err := models.LoadSettings("./omitempty_settings_test.yml")
		So(err, ShouldBeNil)

		Convey("Router section", func() {
			So(models.Set.Router.Login.Path, ShouldEqual, "/login")
			So(models.Set.Router.Register.Path, ShouldEqual, "/register")
			So(models.Set.Router.Validate.Path, ShouldEqual, "/validate")
			So(models.Set.Router.Logout.Path, ShouldEqual, "/logout")
			So(models.Set.Router.Update.Path, ShouldEqual, "/update")
			So(models.Set.Router.ResetPassword.Path, ShouldEqual, "/reset")
		})
	})
}

func TestSslSetting(t *testing.T) {
	models.Set = nil

	Convey("Test valid parse setting file", t, func() {
		So(models.Set, ShouldBeNil)
		err := models.LoadSettings("./settings_test.yml")
		So(err, ShouldBeNil)

		Convey("Ssl section", func() {
			So(models.Set.Ssl.Key, ShouldEqual, "./pass_key")
			So(models.Set.Ssl.Certificate, ShouldEqual, "./pass_certificate")
		})
	})
}

func TestDatabaseSetting(t *testing.T) {
	models.Set = nil

	Convey("Test valid parse setting file", t, func() {
		So(models.Set, ShouldBeNil)
		err := models.LoadSettings("./settings_test.yml")
		So(err, ShouldBeNil)

		Convey("Database section", func() {
			So(models.Set.Database.Host, ShouldEqual, "localhost")
			So(models.Set.Database.DbName, ShouldEqual, "testdb")
			So(models.Set.Database.TokenTable, ShouldEqual, "tokentable")
			So(models.Set.Database.UserTable, ShouldEqual, "usertable")
		})
	})
}

func TestOmitLoadSetting(t *testing.T) {
	models.Set = nil

	Convey("Load settings with empty fields", t, func() {
		So(models.Set, ShouldBeNil)
		err := models.LoadSettings("./omitempty_settings_test.yml")
		So(err, ShouldBeNil)

		Convey("Check all field must be empty", func() {
			So(models.Set.Timber.Token, ShouldEqual, "")
			So(models.Set.Timber.Host, ShouldEqual, "")
			So(models.Set.Ssl.Certificate, ShouldEqual, "")
			So(models.Set.Ssl.Key, ShouldEqual, "")
		})

	})
}
