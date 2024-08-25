package config

import (
	"log"
	"strconv"
	"sync"
)

type (
	Configuration struct {
		Server                     *Server
		Db                         *Db
		CacheDb                    *CacheDb
		ExternalServicesCredential *ExternalServicesCredential
	}
	Server struct {
		Host string
		Port int
	}

	Db struct {
		Host string
		Port int
	}

	CacheDb struct {
		Host     string
		Port     int
		Password string
		Channel  int
	}

	ExternalServicesCredential struct {
		IloIp       string
		IloUsername string
		IloPassword string
	}
)

var (
	once           sync.Once
	configInstance *Configuration
)

func GetConfig() *Configuration {
	once.Do(
		func() {
			env := LoadDotEnv()
			serverPort, err := strconv.Atoi(env.AppPort)
			cacheDbPort, err := strconv.Atoi(env.RedisPort)
			if err != nil {
				log.Fatal("Error parsing env port")
			}
			configInstance = &Configuration{
				Server: &Server{
					Host: env.AppHost,
					Port: serverPort,
				},
				Db: &Db{
					Host: "./database/local.db",
					Port: 0,
				},
				CacheDb: &CacheDb{
					Host:     env.RedisHost,
					Port:     cacheDbPort,
					Password: env.RedisPassword,
					Channel:  0,
				},
				ExternalServicesCredential: &ExternalServicesCredential{
					IloIp:       env.IloIp,
					IloUsername: env.IloUsername,
					IloPassword: env.IloKey,
				},
			}
		},
	)
	return configInstance
}
