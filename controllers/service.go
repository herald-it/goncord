package controllers

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"

	"github.com/herald-it/goncord/models"
	. "github.com/herald-it/goncord/utils"
	"github.com/herald-it/goncord/utils/querying"
)

type ServiceController struct {
	session *mgo.Session
}

func (sc ServiceController) GetDB() *mgo.Database {
	return sc.session.DB(models.Set.Database.DbName)
}

func NewServiceController(s *mgo.Session) *ServiceController {
	return &ServiceController{s}
}

func (sc ServiceController) IsValid(
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params) *HttpError {

	collect := sc.GetDB().C(models.Set.Database.TokenTable)

	if err := r.ParseForm(); err != nil {
		return &HttpError{err, "Post form can not be parsed.", 500}
	}

	token := new(models.DumpToken)
	if err := Fill(token, r.PostForm); err != nil {
		return &HttpError{err, "Post form is not consistent with structure.", 500}
	}

	if token.Token == "" {
		return &HttpError{nil, "Invalid token value.", 500}
	}

	find_dump_token, err := querying.FindDumpToken(token, collect)
	if err != nil || find_dump_token == nil {
		return &HttpError{err, "Token not found.", 500}
	}

	token_pars, err := jwt.Parse(find_dump_token.Token, nil)
	life_time := token_pars.Claims["iat"]

	time_span := time.Now().Unix() - int64(life_time.(float64))
	if time_span > (7 * 24 * 60 * 60) {
		collect.Remove(find_dump_token)
		return &HttpError{nil, "Time token life has expired.", 500}
	}

	usr := new(models.User)
	usr.Id = find_dump_token.UserId

	find_usr, err := querying.FindUser(usr, sc.GetDB().C(models.Set.Database.UserTable))
	if err != nil {
		return &HttpError{err, "User not found.", 500}
	}

	json_usr, err := json.Marshal(find_usr)
	if err != nil {
		return &HttpError{err, "User can not convert to json.", 500}
	}

	w.Write(json_usr)
	return nil
}
