package models

import (
	"github.com/herald-it/goncord/utils"

	"gopkg.in/yaml.v2"

	"io/ioutil"
)

type Setting struct {
	Database struct {
		Host       string
		DbName     string
		TokenTable string
		UserTable  string
	}
}

var Set Setting

func LoadSettings() {
	text, err := ioutil.ReadFile("./settings.yml")
	utils.LogError(err)

	err = yaml.Unmarshal(text, &Set)
	utils.LogError(err)
}
