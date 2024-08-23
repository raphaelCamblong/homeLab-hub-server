package cache

import "github.com/redis/go-redis/v9"

type Database interface {
	GetClient() *redis.Client
}
