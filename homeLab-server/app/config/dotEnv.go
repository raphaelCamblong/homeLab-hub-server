package config

import (
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

type DotEnvConfig struct {
	IloHost     string `mapstructure:"ILO_HOST"`
	IloUsername string `mapstructure:"ILO_USERNAME"`
	IloKey      string `mapstructure:"ILO_KEY"`
	XoApiHost   string `mapstructure:"XO_API_HOST"`
	XoApiKey    string `mapstructure:"XO_API_KEY"`
	RedisHost   string `mapstructure:"REDIS_HOST"`
	RedisKey    string `mapstructure:"REDIS_KEY"`
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
