
package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// CacheRepository handles cache operations.

type CacheRepository struct {
	client *redis.Client
}

// NewCacheRepository creates a new CacheRepository.
func NewCacheRepository(client *redis.Client) *CacheRepository {
	return &CacheRepository{
		client: client,
	}
}

// SetCache sets a value in the cache.
func (r *CacheRepository) SetCache(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

// GetCache gets a value from the cache.
func (r *CacheRepository) GetCache(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil // Key does not exist
	}
	return val, err
}

// DeleteCache deletes a value from the cache.
func (r *CacheRepository) DeleteCache(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
