package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/ilyakaznacheev/cleanenv"
)


var cfg AppConfig
type AppConfig struct {
	App struct {
		Domain string `yaml:"domain" env:"DOMAIN" env-default:"http://localhost:8080/"`
		Port string `yaml:"port" env:"PORT" env-default:"8080"`
		Datacenter int `yaml:"datacenter" env:"DATACENTER" env-default:10`
	}
	Redis struct {
		Endpoint     string `yaml:"endpoint" env:"REDIS_ENDPOINT" env-default:"localhost:6379"`
		Username string `yaml:"username" env:"REDIS_USERNAME" env-default:"default"`
		Password string `yaml:"password" env:"REDIS_PASSWORD" env-default:""`
	}
}

func init() {

	err := cleanenv.ReadConfig("config.yml", &cfg)
	if err != nil {
		log.Error("Error reading configuration")
	}
}

func GetConfig() AppConfig {
	return cfg
}
