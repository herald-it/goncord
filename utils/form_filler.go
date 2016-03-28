package utils

import (
	"github.com/gorilla/schema"
	"net/url"
)

func Fill(i interface{}, form url.Values) {
	decoder := schema.NewDecoder()
	err := decoder.Decode(i, form)
	LogError(err)
}
