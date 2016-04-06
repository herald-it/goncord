package utils_test

import (
	"github.com/herald-it/goncord/utils"
	. "github.com/smartystreets/goconvey/convey"
	"net/url"
	"testing"
)

type TestStruct struct {
	Field  string
	Field2 int
}

func TestFormFiller(t *testing.T) {
	Convey("Test form filler method", t, func() {
		ts := new(TestStruct)
		form := url.Values{}
		form.Set("Field", "form_field")
		form.Set("Field2", "1")
		err := utils.Fill(ts, form)

		Convey("Fill structure from form", func() {
			So(err, ShouldBeNil)
		})

		Convey("Correct parse form", func() {
			So(ts.Field, ShouldEqual, "form_field")
			So(ts.Field2, ShouldEqual, 1)
		})
	})
}
