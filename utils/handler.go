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

func ErrWrap(eh ErrHandler) Handle {
	return Handle(func(w http.ResponseWriter, r *http.Request, p Params) {
		if e := eh(w, r, p); e != nil {
			log.Printf("Error: %v Message: %v Code: %v\n", e.Error, e.Message, e.Code)
			http.Error(w, e.Message, e.Code)
		}
	})
}
