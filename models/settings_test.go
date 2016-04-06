package models_test

import (
	"github.com/herald-it/goncord/models"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestLoadSettings(t *testing.T) {
	Convey("Load settings from test", t, func() {
		err := models.LoadSettings()

		So(err, ShouldNotBeNil)
		So(models.Set.Database, ShouldNotBeNil)
		So(models.Set.Database.DbName, ShouldNotBeNil)
		So(models.Set.Database.Host, ShouldNotBeNil)
		So(models.Set.Database.TokenTable, ShouldNotBeNil)
		So(models.Set.Database.UserTable, ShouldNotBeNil)
	})
}
