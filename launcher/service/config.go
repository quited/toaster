package service

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Service struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		ApiEndpoint string `json:"api_endpoint"`
		ProgramFile string `json:"program_file"`
	} `json:"service"`

	Manager struct {
		ApiEndpoint string `json:"api_endpoint"`
	} `json:"manager"`
}

func LoadConfig(configFile string) (*Config, error) {
	s := Config{}
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &s)
	return &s, err
}

func (c *Config) LoadService() (*Service, error) {
	return NewService(c.Service.Name, c.Manager.ApiEndpoint, c.Service.ApiEndpoint)
}

func (c *Config) InstallService() error {
	return InstallService(c.Service.Name, c.Service.Description, c.Service.ProgramFile)
}

func (c *Config) RemoveService() error {
	return RemoveService(c.Service.Name)
}
