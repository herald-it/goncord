package main

import (
	"flag"
	"log"

	"github.com/herald-it/goncord/controllers"
	"github.com/herald-it/goncord/models"
	. "github.com/herald-it/goncord/utils"

	"net/http"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
)

var (
	settingPath = flag.String("s", "./settings.yml", "setting file path")
)

func getSession() *mgo.Session {
	s, err := mgo.Dial(models.Set.Database.Host)
	if err != nil {
		panic(err)
	}

	return s
}

func init() {
	flag.Parse()

	if err := models.LoadSettings(*settingPath); err != nil {
		panic(err)
	}
}

func main() {
	log.Println("Start auth gate")

	uc := controllers.NewUserController(getSession())
	us := controllers.NewServiceController(getSession())

	var router = httprouter.New()
	router.POST(models.Set.Router.Register, ErrWrap(uc.RegisterUser))
	router.POST(models.Set.Router.Login, ErrWrap(uc.LoginUser))
	router.POST(models.Set.Router.Validate, ErrWrap(us.IsValid))
	router.GET("/", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Write([]byte("Service authorization"))
	})

	if err := http.ListenAndServe(models.Set.IP, router); err != nil {
		panic(err)
	}
}
