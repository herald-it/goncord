package controllers

import (
	"encoding/hex"
	"encoding/json"
	"log"
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

// dumpUser save user and token to table token_dump.
func (uc UserController) dumpUser(usr *models.User, token string) error {
	dumpToken := models.NewDumpToken(usr, token)
	err := uc.GetDB().C(models.Set.Database.TokenTable).Insert(&dumpToken)

	return err
}

// LoginUser user authorization.
// Authorization information is obtained from
// form post. In order to log in
// post the form should contain fields such as:
// 	login
// 	password
// 	email
// If authentication is successful, the user in the cookie
// will add the jwt token. Cook's name will be the jwt and the value
// the issued token.
// The token lifetime is 7 days. After the expiration of
// the lifetime of the token, the authorization process need
// pass again.
func (uc UserController) LoginUser(
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params) *HttpError {

	collect := uc.GetDB().C(models.Set.Database.UserTable)

	if err := r.ParseForm(); err != nil {
		return &HttpError{err, "Post form can not be parsed.", 500}
	}

	usr := new(models.User)
	if err := Fill(usr, r.PostForm, "login|email", "password"); err != nil {
		return &HttpError{err, "Error fill form. Not all fields are specified.", 500}
	}

	usr.Password = hex.EncodeToString(pwd_hash.Sum([]byte(usr.Password)))

	userExist, err := querying.FindUser(usr, collect)
	if userExist == nil || err != nil {
		return &HttpError{err, "User not exist.", 500}
	}

	keyPair, err := keygen.NewKeyPair()
	if err != nil {
		return &HttpError{err, "New key pair error.", 500}
	}

	token, err := userExist.NewToken(keyPair.Private)
	if err != nil {
		return &HttpError{err, "New token error.", 500}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    token,
		Domain:   models.Set.Domain,
		HttpOnly: true,
		Secure:   false}) // TODO: HTTPS. Если true то токена не видно.

	if err = uc.dumpUser(userExist, token); err != nil {
		return &HttpError{err, "Token can not be dumped.", 500}
	}

	log.Println("Token added: ", token)
	usr.Password = usr.Password[:5] + "..."
	log.Println("For user: ", usr)
	return nil
}

// RegisterUser registration of the user.
// Details for registration are obtained from
// form post.
// For registration must be post
// the form contained fields such as:
// 	login
// 	password
// 	email
// After registration the token is not issued.
// To retrieve the token you need to pass the operation
// a login.
func (uc UserController) RegisterUser(
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params) *HttpError {

	collect := uc.GetDB().C(models.Set.Database.UserTable)

	if err := r.ParseForm(); err != nil {
		return &HttpError{err, "Post form can not be parsed.", 500}
	}

	usr := new(models.User)
	if err := Fill(usr, r.PostForm, "login", "email", "password"); err != nil {
		return &HttpError{err, "Error fill form. Not all fields are specified.", 500}
	}

	if usr.Login == "" || usr.Email == "" || usr.Password == "" {
		return &HttpError{nil, "All required fields were not filled.", 500}
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

	usr.Password = usr.Password[:5] + "..."
	log.Println("User added: ", usr)
	return nil
}

func (uc UserController) UpdateUser(
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params) *HttpError {

	collect := uc.GetDB().C(models.Set.Database.UserTable)

	if err := r.ParseForm(); err != nil {
		return &HttpError{err, "Post form can not be parsed.", 500}
	}

	updUsrText := r.PostFormValue("user")
	if updUsrText == "" {
		return &HttpError{nil, "Empty update user form.", 500}
	}

	usr := new(models.User)
	if err := json.Unmarshal([]byte(updUsrText), usr); err != nil {
		return &HttpError{err, "Error unmarshal json to user model.", 500}
	}

	if err := collect.UpdateId(usr.ID, usr); err != nil {
		return &HttpError{err, "Error updating user model.", 500}
	}

	return nil
}
