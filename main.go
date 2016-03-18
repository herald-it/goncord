package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func main() {
	var router = httprouter.New()
	router.GET("/", index)
	log.Fatal(http.ListenAndServe(":8000", router))
}
