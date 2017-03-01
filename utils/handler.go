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
	return Handle(func(writer http.ResponseWriter, request *http.Request, params Params) {
		if e := eh(writer, request, params); e != nil {
			log.Printf("\033[7m\033[1m\t ✗ Error: %v Message: %v Code: %v\033[0m",
				e.Error, e.Message, e.Code)

			if failureURL := request.URL.Query().Get("failure"); failureURL != "" {
				http.Redirect(writer, request, failureURL, 301)
			} else {
				http.Error(writer, e.Message, e.Code)
			}
		} else {
			log.Print("\033[7m\033[1m\t ✓ Successfully.\033[0m")

			if successURL := request.URL.Query().Get("success"); successURL != "" {
				http.Redirect(writer, request, successURL, 301)
			}
		}
	})
}
