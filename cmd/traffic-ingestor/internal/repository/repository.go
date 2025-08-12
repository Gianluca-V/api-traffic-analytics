package repository

import (
	"context"
	"encoding/json"

	"api-traffic-analytics/internal/pkg/postgres"
	"api-traffic-analytics/internal/pkg/redis"
	"api-traffic-analytics/internal/shared/models"
)

type Repository struct {
	db  *postgres.TrafficDataRepository
	rdb *redis.CacheRepository
}

func NewRepository(db *postgres.TrafficDataRepository, rdb *redis.CacheRepository) *Repository {
	return &Repository{db: db, rdb: rdb}
}

func (r *Repository) StoreTrafficData(ctx context.Context, data *models.TrafficData) error {
	// Store in PostgreSQL
	if _, err := r.db.Create(ctx, data); err != nil {
		return err
	}

	// Store in Redis
	key := "latest_traffic:" + data.LocationID
	val, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return r.rdb.SetCache(ctx, key, val, 0)
}
