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
	})
}
