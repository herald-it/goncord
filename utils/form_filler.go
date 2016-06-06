package utils

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

func Fill(i interface{}, form url.Values) error {
	// decoder := schema.NewDecoder()

	loadModel(i, form)
	return nil
}

func loadModel(obj interface{}, m url.Values) {
	defer func() {
		if e := recover(); e != nil {
			fmt.Printf("Panic! %v\n", e)
		}
	}()

	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for k, v := range m {
		if f := val.FieldByName(strings.Title(k)); f.IsValid() {
			if f.CanSet() {

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
			} else {
				fmt.Printf("Key '%s' cannot be set\n", k)
			}
		}
	}
}
