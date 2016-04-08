package main

import (
	"github.com/herald-it/goncord/controllers"
	"github.com/herald-it/goncord/models"
	. "github.com/herald-it/goncord/utils"

	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
)

func getSession() *mgo.Session {
	s, err := mgo.Dial(models.Set.Database.Host)
	if err != nil {
		panic(err)
	}

	return s
}

func init() {
	if err := models.LoadSettings(); err != nil {
		panic(err)
	}
}

func main() {
	uc := controllers.NewUserController(getSession())
	us := controllers.NewServiceController(getSession())

	var router = httprouter.New()
	router.POST(models.Set.Router.Register, ErrWrap(uc.RegisterUser))
	router.POST(models.Set.Router.Login, ErrWrap(uc.LoginUser))
	router.POST(models.Set.Router.Validate, ErrWrap(us.IsValid))

	log.Fatal(http.ListenAndServeTLS(
		":8228",
		models.Set.Ssl.Sertificate,
		models.Set.Ssl.Key,
		router))
}
