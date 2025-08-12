package repository

import (
	"context"

	"api-traffic-analytics/internal/pkg/postgres"
	"api-traffic-analytics/internal/shared/models"
)

type Repository struct {
	postgresRepo *postgres.TrafficDataRepository
}

func NewRepository(postgresRepo *postgres.TrafficDataRepository) *Repository {
	return &Repository{postgresRepo: postgresRepo}
}

func (r *Repository) StoreAnalyticsResult(ctx context.Context, result *models.AnalyticsResult) error {
	return r.postgresRepo.CreateAnalyticsResult(ctx, result)
}
