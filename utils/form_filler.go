package utils

import (
	"net/url"

	"github.com/gorilla/schema"
)

func Fill(i interface{}, form url.Values) error {
	decoder := schema.NewDecoder()

	if err := decoder.Decode(i, form); err != nil {
		return err
	}

	return nil
}

// func decode(dst interface{}, form map[string][]string) error {
// 	v := reflect.ValueOf(dst)
// 	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
// 		return errors.New("schema: interface must be a pointer to struct")
// 	}
// }
