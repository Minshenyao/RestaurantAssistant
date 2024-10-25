package utils

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

type Config struct {
	Database struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		Dbname   string `json:"dbname"`
	} `yaml:"database"`
	Services struct {
		Port string `json:"port"`
	} `yaml:"services"`
}

var AppConfig Config

func LoadConfig() {
	file, err := ioutil.ReadFile("config/config.yaml")
	if err != nil {
		log.Println(err)
	}
	err = yaml.Unmarshal(file, &AppConfig)
	if err != nil {
		log.Println(err)
	}
}
