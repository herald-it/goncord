package controllers

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/herald-it/goncord/models"
	"github.com/herald-it/goncord/pwd_hash"
	"github.com/herald-it/goncord/utils"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
)

type UserController struct {
	session *mgo.Session
}

func NewUserController(s *mgo.Session) *UserController {
	return &UserController{s}
}

func (uc UserController) RegisterUser(
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params) {

	err := r.ParseForm()
	utils.LogError(err)

	usr := new(models.User)
	utils.Fill(usr, r.PostForm)

	usr.Password = hex.EncodeToString(pwd_hash.Sum([]byte(usr.Password)))

	// collect := uc.session.DB("auth_service").C("users")
	// n, err := collect.Find(bson.M{"$or": [...]bson.M{bson.M{"login": login}, bson.M{"email": email}}}).Count()
	// utils.LogError(err)

	// if n != 0 {
	// 	w.Write([]byte("User is already exist!"))
	// 	return
	// }

	// collect.Insert(&usr)
	fmt.Println(usr)
	w.Write([]byte("Succesfully added!"))
}
