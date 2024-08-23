package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
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
}

func LoadDotEnv() *DotEnvConfig {
	path := ".env"
	err := godotenv.Load(path)
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
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
	}
}
