package utils_test

import (
	"github.com/herald-it/goncord/utils"
	"net/url"
	"testing"
)

type TestStruct struct {
	Field  string
	Field2 int
}

func TestFormFiller(t *testing.T) {
	ts := new(TestStruct)
	form := url.Values{}
	form.Set("Field", "form_field")
	form.Set("Field2", "1")

	err := utils.Fill(ts, form)
	if err != nil {
		t.Fatalf("Fill form return err: %v", err.Error())
	}

	if ts.Field != "form_field" {
		t.Fatalf("%v not equal test value: %v", ts.Field, form.Get("Field"))
	}

	if ts.Field2 != 1 {
		t.Fatalf("%v not equal test value: %v", ts.Field2, form.Get("Field2"))
	}
}
