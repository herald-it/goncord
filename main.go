package main

import (
	"flag"
	"log"

	"github.com/herald-it/goncord/controllers"
	"github.com/herald-it/goncord/models"
	. "github.com/herald-it/goncord/utils"

	"net/http"

	"github.com/herald-it/goncord/middleware"
	. "github.com/julienschmidt/httprouter"
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
	log.Println("Start initialize...")

	uc := controllers.NewUserController(getSession())
	us := controllers.NewServiceController(getSession())

	coll := middleware.MidCollect{}
	coll = coll.Add(middleware.CheckPermission).Add(middleware.Logging)

	router := New()
	router.POST(
		models.Set.Router.Register.Path,
		coll.Wrap(ErrWrap(uc.RegisterUser)),
	)
	router.POST(
		models.Set.Router.Login.Path,
		coll.Wrap(ErrWrap(uc.LoginUser)),
	)
	router.POST(
		models.Set.Router.Validate.Path,
		coll.Wrap(ErrWrap(us.IsValid)),
	)
	router.POST(
		models.Set.Router.Logout.Path,
		coll.Wrap(ErrWrap(us.Logout)),
	)
	router.POST(
		models.Set.Router.Update.Path,
		coll.Wrap(ErrWrap(uc.UpdateUser)),
	)

	router.GET("/", func(w http.ResponseWriter, r *http.Request, p Params) {
		w.Write([]byte("Service authorization"))
	})

	log.Println("Start auth gate!")
	if err := http.ListenAndServe(models.Set.IP, router); err != nil {
		panic(err)
	}
}
