package interfaces

import (
	"context"

	"api-traffic-analytics/internal/shared/models"
)

// AnalyticsProcessor define la interfaz para procesamiento de datos
type AnalyticsProcessor interface {
	ProcessTrafficData(ctx context.Context, data *models.TrafficData) error
}
