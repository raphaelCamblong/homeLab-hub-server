package config

import (
	"sync"
)

type (
	Configuration struct {
		App                        *AppConfig
		Db                         *NetworkConnection
		CacheDb                    *NetworkConnection
		ExternalServicesCredential *ExternalServicesCredential
	}

	NetworkConnection struct {
		Host    string
		User    *string
		Key     *string
		Channel *int
	}

	Security struct {
		JwtSecret     string
		JetExpiration int
	}

	ExternalServicesCredential struct {
		Ilo NetworkConnection
		XO  NetworkConnection
	}
)

var (
	once           sync.Once
	configInstance *Configuration
)

func GetConfig() *Configuration {
	once.Do(
		func() {
			env := LoadDotEnv(".env")
			config := LoadYamlConfig("./config.yml")

			helper := int(0)
			configInstance = &Configuration{
				App: config,
				Db: &NetworkConnection{
					Host: "./database/local.db",
				},
				CacheDb: &NetworkConnection{
					Host:    env.RedisHost,
					Key:     &env.RedisKey,
					Channel: &helper,
				},
				ExternalServicesCredential: &ExternalServicesCredential{
					Ilo: NetworkConnection{
						Host: env.IloHost,
						User: &env.IloUsername,
						Key:  &env.IloKey,
					},
					XO: NetworkConnection{
						Host: env.XoApiHost,
						Key:  &env.XoApiKey,
					},
				},
			}
		},
	)
	return configInstance
}
