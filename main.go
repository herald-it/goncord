package main

import (
	"github.com/herald-it/goncord/controllers"
	"github.com/herald-it/goncord/utils"

	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
)

func getSession() *mgo.Session {
	set := utils.GetSettingInstance()
	s, err := mgo.Dial(set.Database["host"])
	utils.LogError(err)

	return s
	return nil
}

func main() {
	uc := controllers.NewUserController(getSession())

	var router = httprouter.New()
	router.POST("/register", uc.RegisterUser)
	log.Fatal(http.ListenAndServe(":8228", router))
}
