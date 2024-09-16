package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

type (
	JWTConfig struct {
		Secret     string `mapstructure:"secret"`
		Expiration int    `mapstructure:"expiration"`
	}

	SecurityConfig struct {
		TrustedProxies []string  `mapstructure:"trustedProxies"`
		Enabled        bool      `mapstructure:"enabled"`
		JWT            JWTConfig `mapstructure:"jwt"`
	}

	AppConfig struct {
		Host     string         `mapstructure:"host"`
		Port     int            `mapstructure:"port"`
		Version  string         `mapstructure:"version"`
		Security SecurityConfig `mapstructure:"security"`
	}
)

func LoadYamlConfig(path string) *AppConfig {
	var config AppConfig

	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./")
	if err := viper.ReadInConfig(); err != nil {
		logrus.Errorf("Error config, %v", err)
		os.Exit(1)
	}
	if err := viper.Unmarshal(&config); err != nil {
		logrus.Errorf("Error config, %v", err)
		os.Exit(1)
	}
	return &config
}
