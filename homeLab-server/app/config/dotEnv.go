package config

import (
	"github.com/sirupsen/logrus"
	"os"

	"github.com/joho/godotenv"
)

type DotEnvConfig struct {
	IloIp         string
	IloUsername   string
	IloKey        string
	RedisPassword string
	RedisHost     string
	RedisPort     string
	AppHost       string
	AppPort       string
	XoApiKey      string
	XoApiHost     string
}

func LoadDotEnv() *DotEnvConfig {
	path := ".env"
	err := godotenv.Load(path)
	if err != nil {
		logrus.Errorf("Error loading .env file: %s", err)
	}
	return &DotEnvConfig{
		IloIp:         os.Getenv("ILO_IP"),
		IloUsername:   os.Getenv("ILO_USERNAME"),
		IloKey:        os.Getenv("ILO_KEY"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisHost:     os.Getenv("REDIS_HOST"),
		RedisPort:     os.Getenv("REDIS_PORT"),
		AppHost:       os.Getenv("APP_HOST"),
		AppPort:       os.Getenv("APP_PORT"),
		XoApiKey:      os.Getenv("XO_API_KEY"),
		XoApiHost:     os.Getenv("XO_API_HOST"),
	}
}
