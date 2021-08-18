package config

/*
	This file and app_config.json govern application behavior.
	Not to be confused with the .env file, which contains information considered secret
*/

import (
	"encoding/json"
	"io/ioutil"
)

type StartParameters struct {
	Address        string
	Port           string
	Mode           []string
	CertificateLoc string
	PrivateKeyLoc  string
	PublicKeyLoc   string
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

// Current modes: CUCKOO_NODIFF, CUCKOO_FIXED_ITEMS, CUCKOO_USE_DATABASE
func GetApplicationMode(conf StartParameters) []string {
	return conf.Mode
}
