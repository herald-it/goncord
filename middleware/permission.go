package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/herald-it/goncord/models"

	. "github.com/julienschmidt/httprouter"
)

// CheckPermission middleware to check
// allowed hosts.
// The list of allowed hosts specify
// in the configuration file.
// A sample configuration file:
//	router:
//  	  register:
//          path: /reg
//	    allowedhost:
//	      - localhost:8000
func CheckPermission(next Handle) Handle {
	return Handle(func(w http.ResponseWriter, r *http.Request, p Params) {
		path := strings.Split(r.RequestURI, "?")[0]
		flagPermission := true

		switch path {
		case models.Set.Router.Register.Path:
			flagPermission = flagPermission && contains(
				r.RemoteAddr,
				models.Set.Router.Register.AllowedHost)
		case models.Set.Router.Login.Path:
			flagPermission = flagPermission && contains(
				r.RemoteAddr,
				models.Set.Router.Login.AllowedHost)
		case models.Set.Router.Validate.Path:
			flagPermission = flagPermission && contains(
				r.RemoteAddr,
				models.Set.Router.Validate.AllowedHost)
		case models.Set.Router.Logout.Path:
			flagPermission = flagPermission && contains(
				r.RemoteAddr,
				models.Set.Router.Logout.AllowedHost)
		}

		if !flagPermission {
			log.Println("Remote address: ", r.RemoteAddr, " request rejected.")
			log.Println("Insufficient permissions.")

			http.Error(w, "Insufficient rights.", 500)
			return
		}

		next(w, r, p)
	})
}

func contains(s string, arr []string) bool {
	if len(arr) == 0 {
		return true
	}

	for _, v := range arr {
		if v == s {
			return true
		}
	}

	return false
}
