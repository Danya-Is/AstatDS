package client

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os/user"
)

type Config struct {
	Endpoints []string
	Cluster   string
}

func path() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	homeDir := usr.HomeDir
	return homeDir + "/AstatConfig.json", nil
}

func (config *Config) Write() error {
	curConfig, err := ReadFromDisk()
	if err != nil {
		log.Println(err)
	}
	if curConfig != nil {
		if config.Cluster == "" {
			config.Cluster = curConfig.Cluster
		}
		if len(config.Endpoints) == 0 {
			config.Endpoints = curConfig.Endpoints
		}
	}

	configData, err := json.Marshal(config)
	if err != nil {
		return err
	}
	configPath, err := path()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(configPath, configData, 0777)
	if err != nil {
		return err
	}
	return nil
}

func ReadFromDisk() (*Config, error) {
	configPath, err := path()
	if err != nil {
		return nil, err
	}
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	data := Config{}
	err = json.Unmarshal([]byte(file), &data)
	return &data, err
}
