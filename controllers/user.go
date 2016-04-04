package controllers

import (
	"encoding/hex"
	"net/http"

	"github.com/herald-it/goncord/keygen"
	"github.com/herald-it/goncord/models"
	"github.com/herald-it/goncord/pwd_hash"
	"github.com/herald-it/goncord/querying"
	"github.com/herald-it/goncord/utils"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
)

type UserController struct {
	session *mgo.Session
}

func (uc UserController) GetDB() *mgo.Database {
	return uc.session.DB("auth_service")
}

func NewUserController(s *mgo.Session) *UserController {
	return &UserController{s}
}

// Save user and token to table token_dump.
func (uc UserController) dumpUser(usr *models.User, token string) error {
	dump_token := models.NewDumpToken(usr, token)
	err := uc.GetDB().C("token_dump").Insert(&dump_token)

	return err
}

func (uc UserController) LoginUser(
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params) {

	collect := uc.GetDB().C("users")

	err := r.ParseForm()
	utils.LogError(err)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		w.Write([]byte("Post form can not be parsed."))
		return
	}

	usr := new(models.User)
	utils.Fill(usr, r.PostForm)

	usr.Password = hex.EncodeToString(pwd_hash.Sum([]byte(usr.Password)))

	user_exist, err := querying.FindUser(usr, collect)
	utils.LogError(err)

	if user_exist == nil {
		http.Error(w, http.StatusText(http.StatusPreconditionFailed), http.StatusPreconditionFailed)
		w.Write([]byte("User does not exist."))
		return
	}

	key_pair, err := keygen.NewKeyPair()
	utils.LogError(err)

	token, err := user_exist.NewToken(key_pair.Private)
	utils.LogError(err)

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    token,
		HttpOnly: true,
		Secure:   true})

	err = uc.dumpUser(user_exist, token)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
		w.Write([]byte("Token can not be dump."))
		return
	}

	w.Write([]byte("Token succesfully added."))
}

func (uc UserController) RegisterUser(
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params) {

	collect := uc.GetDB().C("users")

	err := r.ParseForm()
	utils.LogError(err)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		w.Write([]byte("Post form can not be parsed."))
		return
	}

	usr := new(models.User)
	utils.Fill(usr, r.PostForm)

	usr.Password = hex.EncodeToString(pwd_hash.Sum([]byte(usr.Password)))

	user_exist, err := querying.IsExistUser(usr, collect)
	utils.LogError(err)

	if user_exist {
		http.Error(w, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable)
		w.Write([]byte("User already exist"))
		return
	}

	collect.Insert(&usr)
	w.Write([]byte("Succesfully added"))
}
