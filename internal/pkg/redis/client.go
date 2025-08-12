
package redis

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

var (
	client *redis.Client
	ctx    = context.Background()
)

// GetRedisClient returns a configured Redis client.
func GetRedisClient() (*redis.Client, error) {
	if client != nil {
		return client, nil
	}

	// Retrieve Redis connection details from environment variables
	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB := os.Getenv("REDIS_DB")

	// Convert REDIS_DB to an integer
	db := 0
	if redisDB != "" {
		_, err := fmt.Sscanf(redisDB, "%d", &db)
		if err != nil {
			return nil, fmt.Errorf("invalid REDIS_DB value: %w", err)
		}
	}

	client = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       db,
	})

	// Ping the Redis server to check the connection
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("error connecting to Redis: %w", err)
	}

	return client, nil
}
