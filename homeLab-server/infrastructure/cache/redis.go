package cache

import (
	"context"
	"github.com/sirupsen/logrus"
	"sync"

	"github.com/redis/go-redis/v9"
	"homelab.com/homelab-server/homeLab-server/app/config"
)

type redisDatabase struct {
	Client *redis.Client
}

var (
	once        sync.Once
	redisClient *redisDatabase
)

func NewRedisDatabase() (Database, error) {
	once.Do(
		func() {
			c := config.GetConfig()
			client := redis.NewClient(
				&redis.Options{
					Addr:     c.CacheDb.Host,
					Password: *c.CacheDb.Key,
					DB:       *(c.CacheDb.Channel),
				},
			)

			_, err := client.Ping(context.Background()).Result()
			if err != nil {
				logrus.Errorf("failed to connect to Redis: %w",
					err)
				panic(err)
			}

			logrus.Info("Successfully connected to Redis")
			redisClient = &redisDatabase{Client: client}
		},
	)

	return redisClient, nil
}

func (r *redisDatabase) GetClient() *redis.Client {
	return r.Client
}
