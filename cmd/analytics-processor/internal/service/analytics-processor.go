package service

import (
	"context"
	"fmt"
	"time"

	"api-traffic-analytics/cmd/analytics-processor/internal/repository"
	"api-traffic-analytics/internal/interfaces"
	"api-traffic-analytics/internal/shared"
	"api-traffic-analytics/internal/shared/models"
)

type analyticsProcessor struct {
	repo *repository.Repository
}

func NewAnalyticsProcessor(repo *repository.Repository) interfaces.AnalyticsProcessor {
	return &analyticsProcessor{repo: repo}
}

// ProcessTrafficData procesa datos de tráfico individuales y genera análisis
func (p *analyticsProcessor) ProcessTrafficData(ctx context.Context, data *models.TrafficData) error {
	// Validar datos de entrada
	if err := p.validateTrafficData(data); err != nil {
		return fmt.Errorf("invalid traffic data: %w", err)
	}

	// Realizar análisis
	analyticsResults := p.performAnalysis(data)

	// Almacenar resultados
	for _, result := range analyticsResults {
		if err := p.repo.StoreAnalyticsResult(ctx, result); err != nil {
			return fmt.Errorf("failed to store analytics result: %w", err)
		}
	}

	return nil
}

// validateTrafficData valida los datos de tráfico antes del procesamiento
func (p *analyticsProcessor) validateTrafficData(data *models.TrafficData) error {
	if data == nil {
		return fmt.Errorf("traffic data is nil")
	}

	if data.LocationID == "" {
		return fmt.Errorf("location_id is required")
	}

	if data.VehicleCount < 0 {
		return fmt.Errorf("vehicle_count must be non-negative, got %d", data.VehicleCount)
	}

	if data.AverageSpeed < 0 {
		return fmt.Errorf("average_speed must be non-negative, got %.2f", data.AverageSpeed)
	}

	// Validar timestamp
	if data.Timestamp.IsZero() {
		return fmt.Errorf("timestamp is required")
	}

	// Validar nivel de congestión si está presente
	if data.CongestionLevel != "" {
		validLevels := map[string]bool{
			models.CongestionLow:    true,
			models.CongestionMedium: true,
			models.CongestionHigh:   true,
			models.CongestionSevere: true,
		}
		if !validLevels[data.CongestionLevel] {
			return fmt.Errorf("invalid congestion_level: %s", data.CongestionLevel)
		}
	}

	return nil
}

// performAnalysis realiza múltiples análisis sobre los datos de tráfico
func (p *analyticsProcessor) performAnalysis(data *models.TrafficData) []*models.AnalyticsResult {
	var results []*models.AnalyticsResult

	// Calcular índice de congestión
	congestionIndex := p.calculateCongestionIndex(data)
	results = append(results, congestionIndex)

	// Calcular densidad de tráfico
	trafficDensity := p.calculateTrafficDensity(data)
	results = append(results, trafficDensity)

	// Calcular tasa de flujo
	flowRate := p.calculateFlowRate(data)
	results = append(results, flowRate)

	// Calcular tiempo de viaje
	travelTime := p.calculateTravelTime(data)
	results = append(results, travelTime)

	// Calcular índice de retraso
	delayIndex := p.calculateDelayIndex(data)
	results = append(results, delayIndex)

	return results
}

// calculateCongestionIndex calcula el índice de congestión (0-1)
func (p *analyticsProcessor) calculateCongestionIndex(data *models.TrafficData) *models.AnalyticsResult {
	if data.VehicleCount == 0 {
		return &models.AnalyticsResult{
			AnalysisTimestamp: time.Now(),
			PeriodStart:       data.Timestamp,
			PeriodEnd:         data.Timestamp.Add(5 * time.Minute),
			LocationID:        &data.LocationID,
			MetricType:        models.MetricCongestionIndex,
			Value:             0,
			Unit:              models.UnitIndex,
			ConfidenceLevel:   shared.Float64Ptr(0.95),
		}
	}

	// Normalizar valores basados en umbrales típicos
	// Asumiendo 200 vehículos como congestión máxima esperada
	normalizedCount := float64(data.VehicleCount) / 200.0
	if normalizedCount > 1.0 {
		normalizedCount = 1.0
	}

	// Asumiendo 80 km/h como velocidad libre de flujo
	normalizedSpeed := 0.0
	if data.AverageSpeed > 0 {
		normalizedSpeed = (80.0 - data.AverageSpeed) / 80.0
		if normalizedSpeed < 0 {
			normalizedSpeed = 0
		}
	}

	// Índice de congestión combinado
	congestionIndex := (normalizedCount + normalizedSpeed) / 2.0
	if congestionIndex > 1.0 {
		congestionIndex = 1.0
	}

	return &models.AnalyticsResult{
		AnalysisTimestamp: time.Now(),
		PeriodStart:       data.Timestamp,
		PeriodEnd:         data.Timestamp.Add(5 * time.Minute),
		LocationID:        &data.LocationID,
		MetricType:        models.MetricCongestionIndex,
		Value:             congestionIndex,
		Unit:              models.UnitIndex,
		ConfidenceLevel:   shared.Float64Ptr(0.95),
		Trend:             p.determineTrend(congestionIndex),
	}
}

// calculateTrafficDensity calcula la densidad de tráfico (vehículos/km)
func (p *analyticsProcessor) calculateTrafficDensity(data *models.TrafficData) *models.AnalyticsResult {
	// Asumiendo un segmento de carretera de 1 km para este cálculo
	density := float64(data.VehicleCount) / 1.0 // vehículos/km

	return &models.AnalyticsResult{
		AnalysisTimestamp: time.Now(),
		PeriodStart:       data.Timestamp,
		PeriodEnd:         data.Timestamp.Add(5 * time.Minute),
		LocationID:        &data.LocationID,
		MetricType:        models.MetricTrafficDensity,
		Value:             density,
		Unit:              models.UnitVehiclesPerKm,
		ConfidenceLevel:   shared.Float64Ptr(0.90),
	}
}

// calculateFlowRate calcula la tasa de flujo (vehículos/hora)
func (p *analyticsProcessor) calculateFlowRate(data *models.TrafficData) *models.AnalyticsResult {
	// Asumiendo que los datos representan un intervalo de 5 minutos
	flowRate := float64(data.VehicleCount) * (60.0 / 5.0) // vehículos/hora

	return &models.AnalyticsResult{
		AnalysisTimestamp: time.Now(),
		PeriodStart:       data.Timestamp,
		PeriodEnd:         data.Timestamp.Add(5 * time.Minute),
		LocationID:        &data.LocationID,
		MetricType:        models.MetricFlowRate,
		Value:             flowRate,
		Unit:              models.UnitVehiclesPerHour,
		ConfidenceLevel:   shared.Float64Ptr(0.85),
	}
}

// calculateTravelTime calcula el tiempo estimado de viaje
func (p *analyticsProcessor) calculateTravelTime(data *models.TrafficData) *models.AnalyticsResult {
	// Asumiendo un segmento de 1 km
	var travelTime float64

	if data.AverageSpeed > 0 {
		// Tiempo = Distancia / Velocidad (en horas)
		travelTime = (1.0 / data.AverageSpeed) * 60 // Convertir a minutos
	} else {
		// Si velocidad es 0, asumir tiempo muy alto
		travelTime = 60.0 // 1 hora como valor máximo
	}

	return &models.AnalyticsResult{
		AnalysisTimestamp: time.Now(),
		PeriodStart:       data.Timestamp,
		PeriodEnd:         data.Timestamp.Add(5 * time.Minute),
		LocationID:        &data.LocationID,
		MetricType:        models.MetricTravelTime,
		Value:             travelTime,
		Unit:              models.UnitMinutes,
		ConfidenceLevel:   shared.Float64Ptr(0.80),
	}
}

// calculateDelayIndex calcula el índice de retraso comparado con condiciones ideales
func (p *analyticsProcessor) calculateDelayIndex(data *models.TrafficData) *models.AnalyticsResult {
	// Asumiendo 80 km/h como velocidad ideal
	idealSpeed := 80.0
	currentSpeed := data.AverageSpeed

	var delayIndex float64
	if idealSpeed > 0 && currentSpeed > 0 {
		// Índice de retraso: 0 = sin retraso, 1 = máximo retraso
		delayIndex = (idealSpeed - currentSpeed) / idealSpeed
		if delayIndex < 0 {
			delayIndex = 0
		}
		if delayIndex > 1 {
			delayIndex = 1
		}
	}

	return &models.AnalyticsResult{
		AnalysisTimestamp: time.Now(),
		PeriodStart:       data.Timestamp,
		PeriodEnd:         data.Timestamp.Add(5 * time.Minute),
		LocationID:        &data.LocationID,
		MetricType:        models.MetricDelayIndex,
		Value:             delayIndex,
		Unit:              models.UnitIndex,
		ConfidenceLevel:   shared.Float64Ptr(0.85),
	}
}

// determineTrend determina la tendencia basada en el índice de congestión
func (p *analyticsProcessor) determineTrend(congestionIndex float64) string {
	if congestionIndex > 0.7 {
		return "increasing"
	} else if congestionIndex < 0.3 {
		return "decreasing"
	}
	return "stable"
}
