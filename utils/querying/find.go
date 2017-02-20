package querying

import (
	"errors"

	. "github.com/herald-it/goncord/models"
	"gopkg.in/mgo.v2"
	. "gopkg.in/mgo.v2/bson"
	"log"
)

// FindUserID looking for a user in the collection.
// If the ID was found was more than 1
// user returns an error.
func FindUserID(obj *User, c *mgo.Collection) (*User, error) {
	var results []User

	err := c.FindId(obj.ID).All(&results)
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

// FindUser searches for the user in the collection.
// If found more than 1 user returns an error.
func FindUser(obj *User, c *mgo.Collection) (*User, error) {
	var results []User

	err := c.Find(
		M{
			"$and": []M{
				M{
					"password": obj.Password,
				},
				M{
					"$or": []M{
						M{"login": obj.Login},
						M{"email": obj.Email},
					},
				},
			},
		},
	).All(&results)

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

// FindDumpToken searches the token in the collection.
// If found more than 1 token returns an error.
func FindDumpToken(obj *DumpToken, c *mgo.Collection) (*DumpToken, error) {
	log.Println(obj.Token)

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
