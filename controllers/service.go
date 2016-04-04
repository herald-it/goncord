package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"

	"github.com/herald-it/goncord/models"
	"github.com/herald-it/goncord/querying"
	"github.com/herald-it/goncord/utils"
)

type ServiceController struct {
	session *mgo.Session
}

func (sc ServiceController) GetDB() *mgo.Database {
	return sc.session.DB("auth_service")
}

func NewServiceController(s *mgo.Session) *ServiceController {
	return &ServiceController{s}
}

func (sc ServiceController) IsValid(
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params) {

	collect := sc.GetDB().C("token_dump")

	err := r.ParseForm()
	utils.LogError(err)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		w.Write([]byte("Post form can not be parsed."))
		return
	}

	token := new(models.DumpToken)
	utils.Fill(token, r.PostForm)

	if token.Token == "" {
		w.Write([]byte("Invalid token value."))
		return
	}

	find_dump_token, err := querying.FindDumpToken(token, collect)
	if err != nil {
		w.Write([]byte("Token not found."))
		return
	}

	usr := new(models.User)
	usr.Id = find_dump_token.UserId

	find_usr, err := querying.FindUser(usr, sc.GetDB().C("users"))
	if err != nil {
		w.Write([]byte("User not found."))
		return
	}

	json_usr, err := json.Marshal(find_usr)
	if err != nil {
		w.Write([]byte("User can not convert to json."))
	}

	w.Write(json_usr)
}
