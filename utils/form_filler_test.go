package utils_test

import (
	"net/url"
	"testing"

	"github.com/herald-it/goncord/utils"
	. "github.com/smartystreets/goconvey/convey"
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

		Convey("Test fill object method", func() {
			testStruct := struct {
				Field1 int
				Field2 string
				Field3 float64
			}{}

			err := utils.Fill(
				&testStruct,
				url.Values{
					"Field1": []string{"5"},
				},
				"Field1",
			)

			So(err, ShouldBeNil)
			So(testStruct.Field1, ShouldEqual, 5)

			err = utils.Fill(
				&testStruct,
				url.Values{
					"Field2": []string{"hello"},
					"Field3": []string{"3.14"},
				},
				"Field1|Field2", "Field3",
			)

			So(err, ShouldBeNil)
			So(testStruct.Field1, ShouldEqual, 5)
			So(testStruct.Field2, ShouldEqual, "hello")
			So(testStruct.Field3, ShouldEqual, 3.14)

			err = utils.Fill(
				&testStruct,
				url.Values{
					"Field1": []string{"1"},
					"Field2": []string{"2"},
					"Field3": []string{"3"},
				},
				"Field1|Field2|Field3",
			)

			So(err, ShouldBeNil)
			So(testStruct.Field1, ShouldEqual, 1)
			So(testStruct.Field2, ShouldEqual, "2")
			So(testStruct.Field3, ShouldEqual, 3.0)
		})
	})
}
