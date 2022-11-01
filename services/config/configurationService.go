package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sirupsen/logrus"
)

type AppConfig struct {
	App struct {
		Domain       string `yaml:"domain" env:"DOMAIN" env-default:"http://localhost:8080/"`
		Port         string `yaml:"port" env:"PORT" env-default:"8080"`
		DatacenterId int    `yaml:"datacenterId" env:"DATACENTER_ID" env-default:"10"`
		MachineId    int    `yaml:"machineId" env:"MACHINE_ID" env-default:"-1"`
		PodName      string `yaml:"podName" env:"POD_NAME" env-default:"curto-0"`
	}
	Redis struct {
		Endpoint string `yaml:"endpoint" env:"REDIS_ENDPOINT" env-default:"localhost:6379"`
		Username string `yaml:"username" env:"REDIS_USERNAME" env-default:"default"`
		Password string `yaml:"password" env:"REDIS_PASSWORD" env-default:""`
	}
}

type AppConfigurationService struct {
	configFilePath string
	appConfig      AppConfig
	log            *logrus.Logger
}

func NewAppConfigurationService(filePath string, logger *logrus.Logger) AppConfigurationService {

	return AppConfigurationService{filePath, loadConfig(logger, filePath), logger}
}

func loadConfig(log *logrus.Logger, configFilePath string) AppConfig {

	log.Debugf("Loading configuration from %s", configFilePath)
	var cfg AppConfig
	err := cleanenv.ReadConfig(configFilePath, &cfg)
	if err != nil {
		log.Errorf("Error reading configuration: %s", err.Error())
	}

	return cfg
}

func (a *AppConfigurationService) GetConfig() AppConfig {
	return a.appConfig
}
