package main

import (
	"log"
	"net/http"

	"github.com/herald-it/goncord/controllers"
	"github.com/herald-it/goncord/utils"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
)

func getSession() *mgo.Session {
	set := utils.GetSettingInstance()
	s, err := mgo.Dial(set.Database["host"])
	utils.LogError(err)

	return s
}

func main() {
	uc := controllers.NewUserController(getSession())

	var router = httprouter.New()
	router.POST("/register", uc.RegisterUser)
	log.Fatal(http.ListenAndServe(":8228", router))
}
