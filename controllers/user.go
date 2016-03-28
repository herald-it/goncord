package controllers

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"net/http"

	"github.com/herald-it/goncord/models"
	"github.com/herald-it/goncord/utils"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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

	login := r.PostFormValue("login")
	password := r.PostFormValue("password")
	email := r.PostFormValue("email")

	tmp_u := models.User{}
	collect := uc.session.DB("auth_service").C("users")
	n, err := collect.Find(bson.M{"$or": [...]bson.M{bson.M{"login": login}, bson.M{"email": email}}}).Count()
	utils.LogError(err)

	if n != 0 {
		w.Write([]byte("User is already exist!"))
		return
	}

	hasher_sha1 := sha1.New()
	_, err = hasher_sha1.Write([]byte(password))
	utils.LogError(err)

	hasher_md5 := md5.New()
	_, err = hasher_md5.Write(hasher_sha1.Sum(nil))

	tmp_u = models.User{
		Login:    login,
		Password: hex.EncodeToString(hasher_md5.Sum(nil)),
		Email:    email,
	}

	collect.Insert(&tmp_u)

	w.Write([]byte("Succesfully added!"))
}
