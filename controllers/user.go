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

func (uc *UserController) getUserTable() *mgo.Collection {
	return uc.session.DB(models.Set.Database.DbName).C(models.Set.Database.UserTable)
}

func (uc *UserController) getTokenTable() *mgo.Collection {
	return uc.session.DB(models.Set.Database.DbName).C(models.Set.Database.TokenTable)
}

// NewUserController create new user contgroller.
func NewUserController(s *mgo.Session) *UserController {
	return &UserController{s}
}

// dumpUser save user and token to table token_dump.
func (uc UserController) dumpUser(usr *models.User, token string) error {
	dumpToken := models.NewDumpToken(usr, token)
	err := uc.getTokenTable().Insert(&dumpToken)

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

	collect := uc.getUserTable()

	if err := r.ParseForm(); err != nil {
		return &HttpError{Error: err, Message: "Post form can not be parsed.", Code: 500}
	}

	usr := models.User{}
	if err := Fill(&usr, r.PostForm, "login|email", "password"); err != nil {
		return &HttpError{Error: err, Message: "Error fill form. Not all fields are specified.", Code: 500}
	}

	usr.SetPassword(usr.Password)

	userExist := querying.IsExistUserByLoginOrEmail(
		usr.Login, usr.Email, collect,
	)

	if !userExist {
		return &HttpError{Error: nil, Message: "User does not exist.", Code: 500}
	}

	user, err := querying.FindUser(usr, collect)
	if err != nil {
		return &HttpError{Error: err, Message: "Failed get user model.", Code: 500}
	}

	keyPair, err := keygen.NewKeyPair()
	if err != nil {
		return &HttpError{Error: err, Message: "New key pair error.", Code: 500}
	}

	token, err := user.NewToken(keyPair.Private)
	if err != nil {
		return &HttpError{Error: err, Message: "New token error.", Code: 500}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    token,
		Domain:   models.Set.Domain,
		HttpOnly: true,
		Secure:   false})

	if err = uc.dumpUser(user, token); err != nil {
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

	collect := uc.getUserTable()

	if err := r.ParseForm(); err != nil {
		return &HttpError{Error: err, Message: "Post form can not be parsed.", Code: 500}
	}

	usr := models.User{}
	if err := Fill(&usr, r.PostForm, "login", "email", "password"); err != nil {
		return &HttpError{Error: err, Message: "Error fill form. Not all fields are specified.", Code: 500}
	}

	if usr.Login == "" || usr.Email == "" || usr.Password == "" {
		return &HttpError{Error: nil, Message: "All required fields were not filled.", Code: 500}
	}

	usr.SetPassword(usr.Password)

	isUserExist := querying.IsExistUserByLoginOrEmail(
		usr.Login, usr.Email, collect,
	)

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

	userCollect := uc.getUserTable()

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

	httpError := checkUpdateRules(usr)
	if httpError != nil {
		return httpError
	}

	dumpToken, httpError := uc.getDumpTokenFromRequest(r)
	if httpError != nil {
		return httpError
	}

	usrID := dumpToken.UserId
	if err := userCollect.UpdateId(usrID, bson.M{"$set": usr}); err != nil {
		return &HttpError{Error: err, Message: "Error updating user model.", Code: 500}
	}

	return nil
}

func (uc UserController) ResetPassword(
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params) *HttpError {

	err := r.ParseForm()
	if err != nil {
		return &HttpError{Error: err, Message: "Post form can not be parsed.", Code: 500}
	}

	data := r.PostForm
	oldPassword := data.Get("old_password")
	newPassword := data.Get("new_password")

	user, httpErr := uc.getUserFromRequest(r)
	if httpErr != nil {
		return httpErr
	}

	if oldPassword == "" {
		return uc.forcePasswordChange(user, newPassword)
	} else {
		return uc.passwordChange(user, oldPassword, newPassword)
	}
}

func (uc UserController) forcePasswordChange(user *models.User, password string) *HttpError {
	if password == "" {
		return &HttpError{Error: nil, Message: "Password should be not empty.", Code: 500}
	}

	user.SetPassword(password)
	err := user.Update(uc.getUserTable())
	if err != nil {
		return &HttpError{Error: err, Message: "User update problem.", Code: 500}
	}

	return nil
}

func (uc UserController) passwordChange(user *models.User, oldPassword, newPassword string) *HttpError {
	hashOldPassword := hex.EncodeToString(pwd_hash.Sum([]byte(oldPassword)))

	if hashOldPassword != user.Password {
		return &HttpError{Error: nil, Message: "Old password not equal current password.", Code: 500}
	}

	user.SetPassword(newPassword)
	err := user.Update(uc.getUserTable())
	if err != nil {
		return &HttpError{Error: nil, Message: "User update problem.", Code: 500}
	}

	return nil
}

func checkUpdateRules(usr *models.User) *HttpError {
	if usr.ID != "" || usr.Login != "" || usr.Email != "" {
		return &HttpError{
			Error:   nil,
			Message: "ID, login, email does not update the field.",
			Code:    500}
	}
	if usr.Password != "" {
		return &HttpError{
			Error:   nil,
			Message: "Password does not update field. Please use change password view.",
			Code:    500}
	}

	return nil
}

func (uc UserController) getUserFromRequest(r *http.Request) (*models.User, *HttpError) {
	userCollection := uc.GetDB().C(models.Set.Database.UserTable)

	token, httError := uc.getDumpTokenFromRequest(r)
	if httError != nil {
		return nil, httError
	}

	tmpUsr := &models.User{ID: token.UserId}
	user, err := querying.FindUserID(tmpUsr, userCollection)
	if err != nil {
		return nil, &HttpError{Error: err, Message: "Error find user"}
	}

	return user, nil
}

func (uc UserController) getDumpTokenFromRequest(r *http.Request) (*models.DumpToken, *HttpError) {
	tokenCollect := uc.GetDB().C(models.Set.Database.TokenTable)

	token := &models.DumpToken{}
	tokenTmp, httpErr := getToken(r)
	if httpErr != nil {
		return nil, httpErr
	}
	token.Token = tokenTmp
	findDumpToken, err := querying.FindDumpToken(token, tokenCollect)
	if err != nil || findDumpToken == nil {
		return nil, &HttpError{Error: err, Message: "Token not found.", Code: 500}
	}

	return findDumpToken, nil
}
