package utils

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"github.com/herald-it/goncord/models"
	"net/url"
	"log"
	"gopkg.in/yaml.v2"
)

type ErrorReport struct {
	Error   HttpError
	Request http.Request
	Params  httprouter.Params
}

func TelegramReport(err interface{}) {
	go func() {
		yamlMessage, err := yaml.Marshal(err)

		resp, err := http.PostForm(models.Set.Timber.Host,
			url.Values{
				"token": {models.Set.Timber.Token},
				"message": {string(yamlMessage)},
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
	}()
}