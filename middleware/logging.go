package middleware

import (
	"log"
	"net/http"
	"time"

	. "github.com/julienschmidt/httprouter"
)

func Logging(next Handle) Handle {
	return func(w http.ResponseWriter, r *http.Request, p Params) {
		log.Println("\033[7m\033[1mThe incoming request:\033[0m ", r)

		t0 := time.Now()
		next(w, r, p)
		t1 := time.Now()

		log.Println("The request has been processed.")
		log.Printf("\033[7m\033[1mElapsed time: %v\033[0m", t1.Sub(t0))
	}
}
