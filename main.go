package main

import (
	"github.com/herald-it/goncord/controllers"
	. "github.com/herald-it/goncord/utils"

	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
)

func getSession() *mgo.Session {
	set := GetSettingInstance()
	s, err := mgo.Dial(set.Database.Host)
	LogError(err)

	return s
}

func main() {
	uc := controllers.NewUserController(getSession())
	us := controllers.NewServiceController(getSession())

	var router = httprouter.New()
	router.POST("/register", ErrWrap(uc.RegisterUser))
	router.POST("/login", ErrWrap(uc.LoginUser))
	router.POST("/validate", ErrWrap(us.IsValid))

	log.Fatal(http.ListenAndServe(":8228", router))
}
