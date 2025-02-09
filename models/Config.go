package models

import (
	"gopkg.in/yaml.v3"
	"os"
)

//type Config struct {
//	API   ExApi
//}

var Conf Config

type Config struct {
	Url       string `yaml:"Url"`
	ApiKey    string `yaml:"ApiKey"`
	SecretKey string `yaml:"SecretKey"`
}

func Load() {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		//return nil //, fmt.Errorf("failed to read config file: %v", err)
	}

	//var config Config
	if err := yaml.Unmarshal(data, &Conf); err != nil {
		//return nil //, fmt.Errorf("failed to unmarshal config: %v", err)
	}
	//return &config //, nil
}
