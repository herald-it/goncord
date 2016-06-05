package querying

import (
	. "github.com/herald-it/goncord/models"
	"gopkg.in/mgo.v2"
)

func IsExistDumpToken(obj *DumpToken, c *mgo.Collection) (bool, error) {
	dumpToken, err := FindDumpToken(obj, c)
	return dumpToken != nil && err == nil, err
}

func IsExistUser(obj *User, c *mgo.Collection) (bool, error) {
	usr, err := FindUser(obj, c)
	return usr != nil && err == nil, err
}
