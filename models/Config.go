package models

var Conf Config

type Config struct {
	Url       string `yaml:"Url"`
	ApiKey    string `yaml:"ApiKey"`
	SecretKey string `yaml:"SecretKey"`
}
