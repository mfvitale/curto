package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/ilyakaznacheev/cleanenv"
)


type AppConfig struct {
	App struct {
		Port string `yaml:"port" env:"PORT" env-default:"8080"`
	}
	Redis struct {
		Endpoint     string `yaml:"endpoint" env:"REDIS_ENDPOINT" env-default:"localhost:6379"`
		Username string `yaml:"username" env:"REDIS_USERNAME" env-default:"default"`
		Password string `yaml:"password" env:"REDIS_PASSWORD" env-default:""`
	}
}

var cfg AppConfig
func init() {

	err := cleanenv.ReadConfig("config.yml", &cfg)
	if err != nil {
		log.Error("Error reading configuration")
	}
}

func GetConfig() AppConfig {
	return cfg
}
