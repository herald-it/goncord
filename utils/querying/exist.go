package querying

import (
	. "github.com/herald-it/goncord/models"
	"gopkg.in/mgo.v2"
)

// IsExistDumpToken checks for the existence
// token in the database.
func IsExistDumpToken(obj *DumpToken, c *mgo.Collection) (bool, error) {
	dumpToken, err := FindDumpToken(obj, c)
	return dumpToken != nil && err == nil, err
}

// IsExistUser checks for the existence
// user in the database.
func IsExistUser(obj *User, c *mgo.Collection) (bool, error) {
	usr, err := FindUser(obj, c)
	return usr != nil && err == nil, err
}
