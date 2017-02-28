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

// Logout removes the current token from
// the database. The next validation
// the user is not authorized.
func (sc ServiceController) Logout(
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params) *HttpError {

	collect := sc.GetDB().C(models.Set.Database.TokenTable)
	token := models.DumpToken{}

	tokenTmp, httpErr := getToken(r)
	if httpErr != nil {
		return httpErr
	}
	token.Token = tokenTmp

	if token.Token == "" {
		return &HttpError{Error: nil, Message: "Empty token value.", Code: 500}
	}

	if err := collect.Remove(token); err != nil {
		return &HttpError{Error: err, Message: "Delete token error.", Code: 500}
	}

	return nil
}

// IsValid Check the token for validity.
// The token can be a cookie or transferred
// post the form. First we checked the cookies.
// If the token is valid, the response will contain
// user model in json format.
func (sc ServiceController) IsValid(
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params) *HttpError {

	collect := sc.GetDB().C(models.Set.Database.TokenTable)
	token := &models.DumpToken{}

	tokenTmp, httpErr := getToken(r)
	if httpErr != nil {
		return httpErr
	}
	token.Token = tokenTmp

	if token.Token == "" {
		return &HttpError{Error: nil, Message: "Empty token value.", Code: 500}
	}

	findDumpToken, err := querying.FindDumpToken(token, collect)
	if err != nil || findDumpToken == nil {
		return &HttpError{Error: err, Message: "Token not found.", Code: 500}
	}

	tokenParse, err := jwt.Parse(findDumpToken.Token, nil)
	if checkLifeTime(tokenParse) {
		collect.Remove(findDumpToken)
		return &HttpError{Error: nil, Message: "Time token life has expired.", Code: 500}
	}

	usr := new(models.User)
	usr.ID = findDumpToken.UserId

	findUsr, err := querying.FindUserID(usr, sc.GetDB().C(models.Set.Database.UserTable))
	if err != nil {
		return &HttpError{Error: err, Message: "User not found.", Code: 500}
	}

	jsonUsr, err := json.Marshal(findUsr)
	if err != nil {
		return &HttpError{Error: err, Message: "User can not convert to json.", Code: 500}
	}

	w.Write(jsonUsr)
	return nil
}

// getToken returns the token from the cookie,
// if the cookie is not present in the token, then looking in
// post the form if the token is not exist, then returned
// an empty string and error code.
func getToken(r *http.Request) (string, *HttpError) {
	jwtCookie, err := r.Cookie("jwt")
	if err != nil {
		if err := r.ParseForm(); err != nil {
			return "", &HttpError{Error: err, Message: "Post form can not be parsed.", Code: 500}
		}

		token := r.PostForm.Get("jwt")
		return token, nil
	}

	if jwtCookie.Value == "" {
		return "", &HttpError{Error: nil, Message: "Empty token value.", Code: 500}
	}

	return jwtCookie.Value, nil
}

// checkLifeTime checks the token lifetime.
func checkLifeTime(token *jwt.Token) bool {
	claims := token.Claims.(jwt.MapClaims)

	lifeTime := claims["iat"]
	timeSpan := time.Now().Unix() - int64(lifeTime.(float64))

	return timeSpan > (7 * 24 * 60 * 60)
}
