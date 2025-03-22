package models

var Conf Config

type Config struct {
	Exchanges map[string]ExchangeConf `yaml:"api"`
	Db        Db
}

type ExchangeConf struct {
	Url       string `yaml:"Url"`
	ApiKey    string `yaml:"ApiKey"`
	SecretKey string `yaml:"SecretKey"`
}

type Db struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"dbname"`
	Port     int    `yaml:"port"`
	Path     string `yaml:"search_path"`
	SslMode  string `yaml:"sslmode"`
}
