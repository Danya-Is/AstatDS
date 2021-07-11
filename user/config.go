package user

import (
	"encoding/json"
	"io/ioutil"
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
	configData, err := json.Marshal(config)
	if err != nil {
		return err
	}
	configPath, err := path()
	if err != nil {
		return err
	}
	ioutil.WriteFile(configPath, configData, 0777)
	return nil
}

func Read() (*Config, error) {
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
