package controllers

import (
	//"encoding/base64"
	"encoding/hex"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/herald-it/goncord/models"
	"github.com/herald-it/goncord/pwd_hash"
	"github.com/herald-it/goncord/utils"

	"github.com/dgrijalva/jwt-go"
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

func (uc UserController) LoginUser(
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params) {

	collect := uc.GetDB().C("users")

	err := r.ParseForm()
	utils.LogError(err)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	usr := new(models.User)
	utils.Fill(usr, r.PostForm)

	usr.Password = hex.EncodeToString(pwd_hash.Sum([]byte(usr.Password)))

	n, err := usr.FindWithPwd(collect).Count()
	utils.LogError(err)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if n != 1 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		w.Write([]byte("User not found!"))
		return
	}

	file_pub_key, err := ioutil.ReadFile("public.key")
	utils.LogError(err)
	rsa_pub_key, err := jwt.ParseRSAPublicKeyFromPEM(file_pub_key)
	_ = rsa_pub_key
	utils.LogError(err)

	file_pri_key, err := ioutil.ReadFile("private.key")
	utils.LogError(err)
	rsa_pri_key, err := jwt.ParseRSAPrivateKeyFromPEM(file_pri_key)
	_ = rsa_pri_key
	utils.LogError(err)

	token := jwt.New(jwt.SigningMethodRS256)
	token.Claims["email"] = usr.Email
	token.Claims["login"] = usr.Login
	token.Claims["iat"] = time.Now().Unix()

	log.Print(token.Header)
	log.Print(token.Claims)

	token.Method = jwt.GetSigningMethod("RS256")

	raw_token_strng, _ := token.SigningString()
	log.Print("Raw token: " + raw_token_strng)

	sign, err := token.Method.Sign(raw_token_strng, rsa_pri_key)
	utils.LogError(err)

	log.Print("Sign: " + sign)

	err = token.Method.Verify(strings.Join(strings.Split(raw_token_strng, ".")[0:2], "."), sign, rsa_pub_key)
	utils.LogError(err)
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

	if n == 1 {
		http.Error(w, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable)
		w.Write([]byte("User already exist"))
		return
	}

	collect.Insert(&usr)
	w.Write([]byte("Succesfully added"))
}
