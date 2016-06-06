package models

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Setting struct {
	Database struct {
		Host       string
		DbName     string
		TokenTable string
		UserTable  string
	}
	Ssl struct {
		Key         string
		Certificate string
	}
	Router struct {
		Register string
		Login    string
		Validate string
		Logout   string
	}
	Domain string
	IP     string
}

// Set the set of loaded settings.
var Set Setting

// LoadSettings loads the settings from a file.
func LoadSettings(path string) error {
	text, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	if err = yaml.Unmarshal(text, &Set); err != nil {
		return err
	}

	return nil
}
