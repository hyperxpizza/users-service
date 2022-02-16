package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	Database struct {
		User     string `json:"user"`
		Password string `json:"password"`
		Port     int    `json:"port"`
		Name     string `json:"name"`
		Host     string `json:"host"`
	} `json:"database"`
	UsersService struct {
		Port int    `json:"port"`
		Host string `json:"host"`
	} `json:"usersService"`
}

func NewConfig(pathToFile string) (*Config, error) {
	file, err := os.Open(pathToFile)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var c Config

	err = json.Unmarshal(data, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (c *Config) PrettyPrint() {
	data, _ := json.MarshalIndent(c, "", " ")
	fmt.Println(string(data))
}
