package utils

import (
	"log"
	"net/http"

	. "github.com/julienschmidt/httprouter"
)

type HttpError struct {
	Error   error
	Message string
	Code    int
}

type ErrHandler func(http.ResponseWriter, *http.Request, Params) *HttpError

// ErrWrap processes errors in the query.
// Logging the queries.
// If the url has parameters of success or failure,
// in case of redirect performs on them
// depending on the result of the request.
// Parameters are passed in the request url, for example:
// 	http://localhost:8000/login?success=http://google.com&failure=http://ya.ru
func ErrWrap(eh ErrHandler) Handle {
	return Handle(func(w http.ResponseWriter, r *http.Request, p Params) {
		if e := eh(w, r, p); e != nil {
			log.Printf("Error: %v Message: %v Code: %v\n", e.Error, e.Message, e.Code)

			if failureURL := r.URL.Query().Get("failure"); failureURL != "" {
				http.Redirect(w, r, failureURL, 301)
			} else {
				http.Error(w, e.Message, e.Code)
			}
		} else {
			if successURL := r.URL.Query().Get("success"); successURL != "" {
				http.Redirect(w, r, successURL, 301)
			}
		}
	})
}
