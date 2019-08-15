package main

import (
	"encoding/json"
	"github.com/quited/toaster/launcher/service"
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

func loadConfig() (*Config, error) {
	s := Config{}
	data, err := ioutil.ReadFile(*configFile)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &s)
	return &s, err
}

func (c *Config) LoadService() (*service.Service, error) {
	return service.NewService(c.Service.Name, c.Manager.ApiEndpoint, c.Service.ApiEndpoint)
}

func (c *Config) InstallService() error {
	return service.InstallService(c.Service.Name, c.Service.Description, c.Service.ProgramFile)
}

func (c *Config) RemoveService() error {
	return service.RemoveService(c.Service.Name)
}
