package models

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Router struct {
	Path        string
	AllowedHost []string
}

type Ssl struct {
	Key         string
	Certificate string
}

type RouterPath struct {
	Register      Router
	Login         Router
	Validate      Router
	Logout        Router
	Update        Router
	ResetPassword Router
}

type Timber struct {
	Host  string
	Token string
}

type Database struct {
	Host       string
	DbName     string
	TokenTable string
	UserTable  string
}

type Setting struct {
	Database Database
	Ssl      Ssl
	Router   RouterPath
	Timber   Timber
	Domain   string
	IP       string
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
