package config

import (
	"encoding/json"
	"github.com/ConfWatch/confwatchd/log"
	"io/ioutil"
)

type DatabaseConfig struct {
	Type       string `json:"type"`
	Connection string `json:"connection"`
}

type Configuration struct {
	Dev      bool           `json:"dev"`
	Address  string         `json:"address"`
	Port     int            `json:"port"`
	Hosts    []string       `json:"hosts"`
	Database DatabaseConfig `json:"db"`
}

var Conf = Configuration{Dev: false}

func Load(filename string) error {
	log.Infof("Loading configuration from %s ...", log.Bold(filename))
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(raw, &Conf)
	if err != nil {
		return err
	}

	return nil
}
