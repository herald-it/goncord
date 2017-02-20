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
	"gopkg.in/mgo.v2/bson"
)

// UserController get access for instance mongo db.
type UserController struct {
	session *mgo.Session
}

// GetDB - get current mongo session.
// Return:
// 	current mongo session.
func (uc UserController) GetDB() *mgo.Database {
	return uc.session.DB(models.Set.Database.DbName)
}

// NewUserController create new user contgroller.
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
		return &HttpError{Error: err, Message: "Post form can not be parsed.", Code: 500}
	}

	usr := new(models.User)
	if err := Fill(usr, r.PostForm, "login|email", "password"); err != nil {
		return &HttpError{Error: err, Message: "Error fill form. Not all fields are specified.", Code: 500}
	}

	usr.Password = hex.EncodeToString(pwd_hash.Sum([]byte(usr.Password)))

	userExist, err := querying.FindUser(usr, collect)
	if userExist == nil || err != nil {
		return &HttpError{Error: err, Message: "User does not exist.", Code: 500}
	}

	keyPair, err := keygen.NewKeyPair()
	if err != nil {
		return &HttpError{Error: err, Message: "New key pair error.", Code: 500}
	}

	token, err := userExist.NewToken(keyPair.Private)
	if err != nil {
		return &HttpError{Error: err, Message: "New token error.", Code: 500}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    token,
		Domain:   models.Set.Domain,
		HttpOnly: true,
		Secure:   false})

	if err = uc.dumpUser(userExist, token); err != nil {
		return &HttpError{Error: err, Message: "Token can not be dumped.", Code: 500}
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
		return &HttpError{Error: err, Message: "Post form can not be parsed.", Code: 500}
	}

	usr := new(models.User)
	if err := Fill(usr, r.PostForm, "login", "email", "password"); err != nil {
		return &HttpError{Error: err, Message: "Error fill form. Not all fields are specified.", Code: 500}
	}

	if usr.Login == "" || usr.Email == "" || usr.Password == "" {
		return &HttpError{Error: nil, Message: "All required fields were not filled.", Code: 500}
	}

	usr.Password = hex.EncodeToString(pwd_hash.Sum([]byte(usr.Password)))

	isUserExist, err := querying.IsExistUser(usr, collect)
	if err != nil {
		return &HttpError{Error: err, Message: "Error check user exist.", Code: 500}
	}

	if isUserExist {
		return &HttpError{Error: nil, Message: "User already exist.", Code: 500}
	}

	collect.Insert(&usr)

	usr.Password = usr.Password[:5] + "..."
	log.Println("User added: ", usr)
	return nil
}

// UpdateUser update fields in the user model.
// Update data are taken from form post.
// Form post parameter "user".
// In order that you could update
// model is required _id field.
// Value field is a json user object.
func (uc UserController) UpdateUser(
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params) *HttpError {

	collect := uc.GetDB().C(models.Set.Database.UserTable)

	if err := r.ParseForm(); err != nil {
		return &HttpError{Error: err, Message: "Post form can not be parsed.", Code: 500}
	}

	updUsrText := r.PostFormValue("user")
	if updUsrText == "" {
		return &HttpError{Error: nil, Message: "Empty user field.", Code: 500}
	}

	usr := new(models.User)
	if err := json.Unmarshal([]byte(updUsrText), usr); err != nil {
		return &HttpError{Error: err, Message: "Error unmarshal json to user model.", Code: 500}
	}

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

	usrID := findDumpToken.UserId

	if err := collect.UpdateId(usrID, bson.M{"$set": usr}); err != nil {
		return &HttpError{Error: err, Message: "Error updating user model.", Code: 500}
	}

	return nil
}
