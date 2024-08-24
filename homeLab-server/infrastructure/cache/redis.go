package cache

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/redis/go-redis/v9"
	"homelab.com/homelab-server/homeLab-server/app/config"
)

type redisDatabase struct {
	client *redis.Client
}

var (
	once        sync.Once
	redisClient *redisDatabase
)

func NewRedisDatabase() (Database, error) {
	once.Do(
		func() {
			c := config.GetConfig()
			addr := fmt.Sprintf("%s:%d", c.CacheDb.Host, c.CacheDb.Port)
			client := redis.NewClient(
				&redis.Options{
					Addr:     addr,
					Password: c.CacheDb.Password,
					DB:       c.CacheDb.Channel,
				},
			)

			_, err := client.Ping(context.Background()).Result()
			if err != nil {
				log.Fatal(
					fmt.Errorf(
						"failed to connect to Redis: %w",
						err,
					),
				) // Change to proper error handling in production
			}

			log.Print("Successfully connected to Redis")
			redisClient = &redisDatabase{client: client}
		},
	)

	return redisClient, nil
}

func (r *redisDatabase) GetClient() *redis.Client {
	return r.client
}
