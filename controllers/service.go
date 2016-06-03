package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"

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
	token := &models.DumpToken{}

	jwtCookie, cookieErr := r.Cookie("jwt")
	if cookieErr != nil {
		if err := r.ParseForm(); err != nil {
			return &HttpError{err, "Post form can not be parsed.", 500}
		}

		if err := Fill(token, r.PostForm); err != nil {
			return &HttpError{err, "Post form is not consistent with structure.", 500}
		}
	} else {
		token.Token = jwtCookie.Value
	}

	if token.Token == "" {
		return &HttpError{nil, "Invalid token value.", 500}
	}

	findDumpToken, err := querying.FindDumpToken(token, collect)
	if err != nil || findDumpToken == nil {
		return &HttpError{err, "Token not found.", 500}
	}

	tokenPars, err := jwt.Parse(findDumpToken.Token, nil)
	lifeTime := tokenPars.Claims["iat"]

	timeSpan := time.Now().Unix() - int64(lifeTime.(float64))
	if timeSpan > (7 * 24 * 60 * 60) {
		collect.Remove(findDumpToken)
		return &HttpError{nil, "Time token life has expired.", 500}
	}

	usr := new(models.User)
	usr.Id = findDumpToken.UserId

	findUsr, err := querying.FindUser(usr, sc.GetDB().C(models.Set.Database.UserTable))
	if err != nil {
		return &HttpError{err, "User not found.", 500}
	}

	jsonUsr, err := json.Marshal(findUsr)
	if err != nil {
		return &HttpError{err, "User can not convert to json.", 500}
	}

	w.Write(jsonUsr)
	return nil
}
