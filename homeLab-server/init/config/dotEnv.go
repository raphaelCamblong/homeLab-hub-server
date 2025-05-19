package config

import (
	"os"

	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type DotEnvConfig struct {
	RedisKey   string `mapstructure:"REDIS_KEY"`
	DbHost     string `mapstructure:"DB_HOST"`
	DbPort     string `mapstructure:"DB_PORT"`
	DbUser     string `mapstructure:"DB_USER"`
	DbPassword string `mapstructure:"DB_PASSWORD"`
	DbName     string `mapstructure:"DB_NAME"`
}

func LoadDotEnv(path string) *DotEnvConfig {
	var config DotEnvConfig
	var result map[string]interface{}
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		logrus.Errorf("Error reading config file, %s", err)
		os.Exit(1)
	}
	err := viper.Unmarshal(&result)
	if err != nil {
		logrus.Errorf("Unable to decode into map, %v", err)
		os.Exit(1)
	}

	decErr := mapstructure.Decode(result, &config)

	if decErr != nil {
		logrus.Errorf("error decoding")
		os.Exit(1)
	}

	return &config
}
