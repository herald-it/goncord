package querying

import (
	"errors"
	. "github.com/herald-it/goncord/models"
	"gopkg.in/mgo.v2"
)

func FindUser(obj *User, c *mgo.Collection) (*User, error) {
	var results []User

	err := c.Find(obj).All(&results)
	if err != nil {
		return nil, err
	}

	if len(results) > 1 {
		return nil, errors.New("Find user returned more 1 record.")
	}

	if len(results) == 0 {
		return nil, nil
	}

	return &results[0], nil
}

func FindDumpToken(obj *DumpToken, c *mgo.Collection) (*DumpToken, error) {
	var results []DumpToken

	err := c.Find(obj).All(&results)
	if err != nil {
		return nil, err
	}

	if len(results) > 1 {
		return nil, errors.New("Find user returned more 1 record.")
	}

	if len(results) == 0 {
		return nil, nil
	}

	return &results[0], nil
}

func IsExistDumpToken(obj *DumpToken, c *mgo.Collection) (bool, error) {
	dump_token, err := FindDumpToken(obj, c)
	return dump_token != nil && err == nil, err
}

func IsExistUser(obj *User, c *mgo.Collection) (bool, error) {
	usr, err := FindUser(obj, c)
	return usr != nil && err == nil, err
}
