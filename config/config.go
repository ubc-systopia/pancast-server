package config

import (
	"encoding/json"
	"io/ioutil"
)

type StartParameters struct {
	Address        string
	Port           string
	CertificateLoc string
	PrivateKeyLoc  string
}

func ReadJSONConfig(filename string, config interface{}) error {
	configData, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(configData, config)
	if err != nil {
		return err
	}
	return nil
}

func GetServerURL(conf StartParameters) string {
	return conf.Address + ":" + conf.Port
}