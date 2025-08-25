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

	// Try to load from .env file first
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err == nil {
		var result map[string]interface{}
		err := viper.Unmarshal(&result)
		if err == nil {
			if decErr := mapstructure.Decode(result, &config); decErr == nil {
				logrus.Info("Loaded configuration from .env file")
			}
		}
	}

	// Environment variables take precedence over .env file values
	// This allows Kubernetes secrets to override .env file values
	if envVal := os.Getenv("REDIS_KEY"); envVal != "" {
		config.RedisKey = envVal
	}
	if envVal := os.Getenv("DB_HOST"); envVal != "" {
		config.DbHost = envVal
	}
	if envVal := os.Getenv("DB_PORT"); envVal != "" {
		config.DbPort = envVal
	}
	if envVal := os.Getenv("DB_USER"); envVal != "" {
		config.DbUser = envVal
	}
	if envVal := os.Getenv("DB_PASSWORD"); envVal != "" {
		config.DbPassword = envVal
	}
	if envVal := os.Getenv("DB_NAME"); envVal != "" {
		config.DbName = envVal
	}

	logrus.Info("Final configuration loaded (environment variables take precedence)")
	logrus.Debugf("Final config values: DB_HOST=%s, DB_PORT=%s, DB_USER=%s, DB_NAME=%s, DB_PASSWORD=%s",
		config.DbHost, config.DbPort, config.DbUser, config.DbName,
		func() string {
			if config.DbPassword != "" {
				return "[SET]"
			} else {
				return "[NOT SET]"
			}
		}())
	return &config
}
