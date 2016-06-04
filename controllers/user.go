package controllers

import (
	"encoding/hex"
	"net/http"

	"github.com/herald-it/goncord/models"
	. "github.com/herald-it/goncord/utils"
	"github.com/herald-it/goncord/utils/keygen"
	"github.com/herald-it/goncord/utils/pwd_hash"
	"github.com/herald-it/goncord/utils/querying"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
)

type UserController struct {
	session *mgo.Session
}

func (uc UserController) GetDB() *mgo.Database {
	return uc.session.DB(models.Set.Database.DbName)
}

func NewUserController(s *mgo.Session) *UserController {
	return &UserController{s}
}

// Save user and token to table token_dump.
func (uc UserController) dumpUser(usr *models.User, token string) error {
	dump_token := models.NewDumpToken(usr, token)
	err := uc.GetDB().C(models.Set.Database.TokenTable).Insert(&dump_token)

	return err
}

func (uc UserController) LoginUser(
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params) *HttpError {

	collect := uc.GetDB().C(models.Set.Database.UserTable)

	if err := r.ParseForm(); err != nil {
		return &HttpError{err, "Post form can not be parsed.", 500}
	}

	usr := new(models.User)
	if err := Fill(usr, r.PostForm); err != nil {
		return &HttpError{err, "Post form is not consistent with structure.", 500}
	}

	usr.Password = hex.EncodeToString(pwd_hash.Sum([]byte(usr.Password)))

	userExist, err := querying.FindUser(usr, collect)

	if userExist == nil || err != nil {
		return &HttpError{err, "User not exist.", 500}
	}

	key_pair, err := keygen.NewKeyPair()
	if err != nil {
		return &HttpError{err, "New key pair error.", 500}
	}

	token, err := userExist.NewToken(key_pair.Private)
	if err != nil {
		return &HttpError{err, "New token error.", 500}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    token,
		Domain:   models.Set.Domain,
		HttpOnly: false,
		Secure:   false}) // TODO: HTTPS. Если true то токена не видно.

	if err = uc.dumpUser(userExist, token); err != nil {
		return &HttpError{err, "Token can not be dumped.", 500}
	}

	w.Write([]byte("Token succesfully added."))
	return nil
}

func (uc UserController) RegisterUser(
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params) *HttpError {

	collect := uc.GetDB().C(models.Set.Database.UserTable)

	if err := r.ParseForm(); err != nil {
		return &HttpError{err, "Post form can not be parsed.", 500}
	}

	usr := new(models.User)
	if err := Fill(usr, r.PostForm); err != nil {
		return &HttpError{err, "Post form is not consistent with structure.", 500}
	}

	usr.Password = hex.EncodeToString(pwd_hash.Sum([]byte(usr.Password)))

	isUserExist, err := querying.IsExistUser(usr, collect)
	if err != nil {
		return &HttpError{err, "Error check user exist.", 500}
	}

	if isUserExist {
		return &HttpError{nil, "User already exist.", 500}
	}

	collect.Insert(&usr)
	w.Write([]byte("Succesfully added"))
	return nil
}
