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

func LoadSettings() error {
	text, err := ioutil.ReadFile("./settings.yml")
	if err != nil {
		return err
	}

	if err = yaml.Unmarshal(text, &Set); err != nil {
		return err
	}

	return nil
}
