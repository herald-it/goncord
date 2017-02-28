package utils

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"github.com/davecgh/go-spew/spew"
	"github.com/herald-it/goncord/models"
	"net/url"
	"log"
)

type ErrorReport struct {
	Error   HttpError
	Request http.Request
	Params  httprouter.Params
}

func telegramReport(error ErrorReport) {
	prettyText := spew.Sdump(error)

	resp, err := http.PostForm(models.Set.Timber.Host,
		url.Values{
			"token": models.Set.Timber.Token,
			"message": prettyText,
		},
	)

	if err != nil {
		log.Println(err)
	}

	if resp.StatusCode == 400 {
		log.Println("Bad token.")
	}

	if resp.StatusCode == 500 {
		log.Println("Timber bot error.")
	}

	if resp.StatusCode == 200 {
		log.Println("Error succesfully sended.")
	}
}