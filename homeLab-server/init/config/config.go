package config

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
)

type API struct {
	Host     string         `toml:"host"`
	Port     int            `toml:"port"`
	Version  string         `toml:"version"`
	Security SecurityConfig `toml:"security"`
}

type JWTConfig struct {
	Secret     string `toml:"secret"`
	Expiration int    `toml:"expiration"`
}

type SecurityConfig struct {
	TrustedProxies []string  `toml:"trustedProxies"`
	Enabled        bool      `toml:"enabled"`
	JWT            JWTConfig `toml:"jwt"`
}

type CoreConfig struct {
	API API `toml:"api"`
}

type LogConfig struct {
	Level string `toml:"level"`
	File  string `toml:"file"`
}

type DatabaseConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	DBName   string `toml:"db_name"`
}

type LambdaConfig struct {
	BaseURL string `toml:"base_url"`
}

type K8sConfig struct {
	Kubeconfig string `toml:"kubeconfig"`
}

type ClientConfig struct {
	Lambda   LambdaConfig   `toml:"lambda"`
	Database DatabaseConfig `toml:"database"`
	K8s      K8sConfig      `toml:"k8s"`
}

type Config struct {
	Core   CoreConfig   `toml:"core"`
	Client ClientConfig `toml:"client"`
	Log    LogConfig    `toml:"log"`
}

func DefaultConfig() *Config {
	return &Config{
		Log: LogConfig{
			Level: "info",
			File:  "/var/log/homelab-server/app.log",
		},
	}
}

func (c *Config) FlattenConfig(dotEnvConfig DotEnvConfig) error {
	logrus.Debugf("TOML Database config: Host=%s, Port=%d, User=%s, DBName=%s",
		c.Client.Database.Host, c.Client.Database.Port, c.Client.Database.User, c.Client.Database.DBName)
	logrus.Debugf("Environment variables: DB_HOST=%s, DB_PORT=%s, DB_USER=%s, DB_NAME=%s",
		dotEnvConfig.DbHost, dotEnvConfig.DbPort, dotEnvConfig.DbUser, dotEnvConfig.DbName)

	// Only override TOML values if environment variables are set
	if dotEnvConfig.DbHost != "" {
		c.Client.Database.Host = dotEnvConfig.DbHost
		logrus.Debug("Overriding Host with environment variable")
	}
	if dotEnvConfig.DbPort != "" {
		port, err := strconv.Atoi(dotEnvConfig.DbPort)
		if err != nil {
			return fmt.Errorf("❌ Failed to convert db port to int: %w", err)
		}
		c.Client.Database.Port = port
		logrus.Debug("Overriding Port with environment variable")
	}
	if dotEnvConfig.DbUser != "" {
		c.Client.Database.User = dotEnvConfig.DbUser
		logrus.Debug("Overriding User with environment variable")
	}
	if dotEnvConfig.DbPassword != "" {
		c.Client.Database.Password = dotEnvConfig.DbPassword
		logrus.Debug("Overriding Password with environment variable")
	}
	if dotEnvConfig.DbName != "" {
		c.Client.Database.DBName = dotEnvConfig.DbName
		logrus.Debug("Overriding DBName with environment variable")
	}

	logrus.Debugf("Final Database config: Host=%s, Port=%d, User=%s, DBName=%s",
		c.Client.Database.Host, c.Client.Database.Port, c.Client.Database.User, c.Client.Database.DBName)
	return nil
}

func (c *Config) FlattenConfigQuick() error {
	dotEnvConfig := LoadDotEnv(".env")
	if err := c.FlattenConfig(*dotEnvConfig); err != nil {
		logrus.Error("❌ Failed to flatten config: %w", err)
		panic(err)
	}
	return nil
}

var (
	instance *Config
	once     sync.Once
)

func Get() *Config {
	once.Do(func() {
		args := ParseArgs()

		cfg, err := LoadConfig(args.ConfigPath)
		if err != nil {
			panic(fmt.Errorf("❌ config load failed: %w", err))
		}
		err = cfg.Validate()
		if err != nil {
			panic(fmt.Errorf("❌ config validation failed: %w", err))
		}
		instance = cfg
	})
	return instance
}

func LoadConfig(path string) (*Config, error) {
	var cfg Config

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to open config file: %w", err)
	}
	defer f.Close()

	if _, err := toml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("❌ Failed to decode config: %w", err)
	}

	return &cfg, nil
}

func (c *Config) Validate() error {
	return nil
}
