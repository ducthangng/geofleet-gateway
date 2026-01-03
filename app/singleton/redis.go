package singleton

import (
	"context"
	"fmt"
	"sync"
	"time"

	redis "github.com/go-redis/redis/v8"
)

var (
	redisClient *redis.Client
	redisOnce   sync.Once
)

// GetRedisClient returns a singleton instance of the Redis Client
func GetRedisClient() (*redis.Client, error) {
	var err error

	redisOnce.Do(func() {
		// 1. Initialize the client
		cfg := GetGlobalConfig()
		client := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
			Username: cfg.RedisUser,
			Password: cfg.RedisPass,
			DB:       cfg.RedisDB,
		})

		// 2. Test the connection with a timeout context
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if _, pingErr := client.Ping(ctx).Result(); pingErr != nil {
			err = fmt.Errorf("failed to connect to redis: %v", pingErr)
			return
		}

		redisClient = client
	})

	return redisClient, err
}
