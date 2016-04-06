package utils

import (
	"github.com/gorilla/schema"
	"net/url"
)

func Fill(i interface{}, form url.Values) error {
	decoder := schema.NewDecoder()
	if err := decoder.Decode(i, form); err != nil {
		return err
	}

	return nil
}
