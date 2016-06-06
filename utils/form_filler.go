package utils

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// Fill fills the structure with data from a form post.
// obj - the object to be filled.
// m - post form.
// required - the field which should be in post form, in
// if they are not, the structure is not filled
// and return the error "Required fields not found.".
func Fill(obj interface{}, m url.Values, required ...string) error {
	defer func() {
		if e := recover(); e != nil {
			panic(e)
		}
	}()

	for _, reqField := range required {
		_, exist := m[reqField]
		if !exist {
			return errors.New("Required fields not found.")
		}
	}

	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for k, v := range m {
		var f reflect.Value

		if f = val.FieldByName(strings.Title(k)); !f.IsValid() {
			continue
		}

		if !f.CanSet() {
			fmt.Printf("Key '%s' cannot be set\n", k)
			continue
		}

		switch f.Type().Kind() {
		case reflect.Int:
			if i, e := strconv.ParseInt(v[0], 0, 0); e == nil {
				f.SetInt(i)
			} else {
				fmt.Printf("Could not set int value of %s: %s\n", k, e)
			}
		case reflect.Float64:
			if fl, e := strconv.ParseFloat(v[0], 0); e == nil {
				f.SetFloat(fl)
			} else {
				fmt.Printf("Could not set float64 value of %s: %s\n", k, e)
			}
		case reflect.String:
			f.SetString(v[0])

		default:
			fmt.Printf("Unsupported format %v for field %s\n", f.Type().Kind(), k)
		}
	}

	return nil
}
