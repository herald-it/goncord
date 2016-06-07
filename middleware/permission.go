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
//	      - localhost
func CheckPermission(next Handle) Handle {
	return func(w http.ResponseWriter, r *http.Request, p Params) {
		path := strings.Split(r.RequestURI, "?")[0]
		flagPermission := true

		remoteAddr := getIP(r.RemoteAddr)

		switch path {
		case models.Set.Router.Register.Path:
			flagPermission = flagPermission && contains(
				remoteAddr,
				models.Set.Router.Register.AllowedHost)
		case models.Set.Router.Login.Path:
			flagPermission = flagPermission && contains(
				remoteAddr,
				models.Set.Router.Login.AllowedHost)
		case models.Set.Router.Validate.Path:
			flagPermission = flagPermission && contains(
				remoteAddr,
				models.Set.Router.Validate.AllowedHost)
		case models.Set.Router.Logout.Path:
			flagPermission = flagPermission && contains(
				remoteAddr,
				models.Set.Router.Logout.AllowedHost)
		}

		if !flagPermission {
			log.Println("Remote address: ", remoteAddr, " request rejected.")
			log.Println("Insufficient permissions.")

			http.Error(w, "Insufficient rights.", 500)
			return
		}

		next(w, r, p)
	}
}

func getIP(addr string) string {
	tmp := strings.Split(addr, ":")

	if len(tmp) >= 1 {
		return tmp[0]
	}

	return "Invalid remote address."
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
