package utils

import (
	"github.com/herald-it/goncord/models"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// Get setting map, from settings.yml file
func GetSettingInstance() models.Setting {
	text, err := ioutil.ReadFile("./settings.yml")
	LogError(err)

	var set models.Setting

	err = yaml.Unmarshal(text, &set)
	LogError(err)

	return set
}
