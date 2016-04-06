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
	LogError(err)

	return s
}

func main() {
	if err := models.LoadSettings(); err != nil {
		panic(err)
	}

	uc := controllers.NewUserController(getSession())
	us := controllers.NewServiceController(getSession())

	var router = httprouter.New()
	router.POST("/register", ErrWrap(uc.RegisterUser))
	router.POST("/login", ErrWrap(uc.LoginUser))
	router.POST("/validate", ErrWrap(us.IsValid))

	log.Fatal(http.ListenAndServe(":8228", router))
}
