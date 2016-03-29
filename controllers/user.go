package controllers

import (
	"encoding/hex"
	"net/http"

	"github.com/herald-it/goncord/models"
	"github.com/herald-it/goncord/pwd_hash"
	"github.com/herald-it/goncord/utils"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
)

const collect = uc.GetDB().C("users")

type UserController struct {
	session *mgo.Session
}

func (uc UserController) GetDB() *mgo.Database {
	return uc.session.DB("auth_service")
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

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	usr := new(models.User)
	utils.Fill(usr, r.PostForm)

	usr.Password = hex.EncodeToString(pwd_hash.Sum([]byte(usr.Password)))

	n, err := usr.Find(collect).Count()
	utils.LogError(err)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if n != 0 {
		http.Error(w, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable)
		w.Write([]byte("User already exist"))
		return
	}

	collect.Insert(&usr)
	w.Write([]byte("Succesfully added"))
}
