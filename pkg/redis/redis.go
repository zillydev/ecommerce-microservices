package redis

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func NewRedisClient(redisURL string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: "",
		DB:       0,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	return rdb
}
