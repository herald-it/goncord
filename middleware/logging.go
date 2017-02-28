package middleware

import (
	"log"
	"net/http"
	"time"

	. "github.com/julienschmidt/httprouter"
)

// Logging middleware produces logging of queries.
// Prints the request processing time.
func Logging(next Handle) Handle {
	return func(w http.ResponseWriter, r *http.Request, p Params) {
		log.Println("\033[7m\033[1mIncoming request.\033[0m")
		log.Printf("\t\033[7m\033[1mMethod:\033[0m %v", r.Method)
		log.Printf("\t\033[7m\033[1mForm:\033[0m %v", r.Form)
		log.Printf("\t\033[7m\033[1mHost:\033[0m %v", r.Host)
		log.Printf("\t\033[7m\033[1mRemote address:\033[0m %v", r.RemoteAddr)
		log.Printf("\t\033[7m\033[1mRequest URI:\033[0m %v", r.RequestURI)
		log.Printf("\t\033[7m\033[1mUser agent:\033[0m %v", r.UserAgent())
		log.Printf("\t\033[7m\033[1mHeaders:\033[0m %v", r.Header)

		t0 := time.Now()
		next(w, r, p)
		t1 := time.Now()

		log.Printf("\t\033[7m\033[1m ⌛ Elapsed time: %v\033[0m", t1.Sub(t0))
	}
}
