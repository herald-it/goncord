package models_test

import (
	"encoding/json"
	"github.com/herald-it/goncord/models"
	"testing"
)

func TestNewDumpTokenModel(t *testing.T) {
	dump_token := &models.DumpToken{}

	if dump_token == nil {
		t.Fatal("Nil pointer after create new dump token")
	}
}

func TestJsonDumpTokenModel(t *testing.T) {
	dump_token := models.DumpToken{
		Token: "my_secret_token"}

	const str = `{"token":"my_secret_token"}`
	b, e := json.Marshal(&dump_token)

	if e != nil {
		t.Fatalf("Error: %v", e.Error())
	}

	if string(b) != str {
		t.Fatalf("%v not equal %v", string(b), str)
	}
}
