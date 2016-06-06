package querying

import (
	"errors"

	. "github.com/herald-it/goncord/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func FindUserID(obj *User, c *mgo.Collection) (*User, error) {
	var results []User

	err := c.FindId(obj.Id).All(&results)
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

func FindUser(obj *User, c *mgo.Collection) (*User, error) {
	var results []User

	err := c.Find(
		bson.M{"$or": []bson.M{
			bson.M{"login": obj.Login},
			bson.M{"email": obj.Email},
		}}).All(&results)
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
