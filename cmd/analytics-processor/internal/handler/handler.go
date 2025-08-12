package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"api-traffic-analytics/internal/interfaces"
	"api-traffic-analytics/internal/shared/models"
)

// MessageHandler se encarga de procesar mensajes individuales de Kafka
type MessageHandler struct {
	analyticsService interfaces.AnalyticsProcessor // Cambiado a la interfaz
	logger           *log.Logger
}

// NewMessageHandler crea una nueva instancia de MessageHandler
func NewMessageHandler(analyticsService interfaces.AnalyticsProcessor, logger *log.Logger) *MessageHandler {
	if logger == nil {
		logger = log.Default()
	}
	return &MessageHandler{
		analyticsService: analyticsService,
		logger:           logger,
	}
}

// HandleMessage procesa un mensaje individual de Kafka
func (h *MessageHandler) HandleMessage(ctx context.Context, msgValue []byte) error {
	startTime := time.Now()

	// Parsear datos de tráfico
	trafficData, err := h.parseTrafficData(msgValue)
	if err != nil {
		return fmt.Errorf("failed to parse traffic data: %w", err)
	}

	// Log inicio de procesamiento
	h.logger.Printf("Processing traffic data for location %s at %v",
		trafficData.LocationID, trafficData.Timestamp)

	// Procesar con el servicio de análisis
	if err := h.analyticsService.ProcessTrafficData(ctx, trafficData); err != nil {
		h.logger.Printf("Error processing analytics for location %s: %v",
			trafficData.LocationID, err)
		return fmt.Errorf("analytics processing failed: %w", err)
	}

	// Log éxito
	duration := time.Since(startTime)
	h.logger.Printf("Successfully processed analytics for location %s in %v",
		trafficData.LocationID, duration)

	return nil
}

// parseTrafficData convierte bytes a modelo de TrafficData
func (h *MessageHandler) parseTrafficData(data []byte) (*models.TrafficData, error) {
	var trafficData models.TrafficData

	if err := json.Unmarshal(data, &trafficData); err != nil {
		return nil, fmt.Errorf("JSON unmarshal failed: %w", err)
	}

	// Validaciones adicionales
	if err := h.validateTrafficData(&trafficData); err != nil {
		return nil, fmt.Errorf("data validation failed: %w", err)
	}

	return &trafficData, nil
}

// validateTrafficData valida los datos de tráfico
func (h *MessageHandler) validateTrafficData(data *models.TrafficData) error {
	if data.LocationID == "" {
		return fmt.Errorf("location_id is required")
	}

	if data.VehicleCount < 0 {
		return fmt.Errorf("vehicle_count must be non-negative, got %d", data.VehicleCount)
	}

	if data.AverageSpeed < 0 {
		return fmt.Errorf("average_speed must be non-negative, got %.2f", data.AverageSpeed)
	}

	// Validar nivel de congestión si está presente
	if data.CongestionLevel != "" {
		validLevels := map[string]bool{
			"low":    true,
			"medium": true,
			"high":   true,
			"severe": true,
		}
		if !validLevels[data.CongestionLevel] {
			return fmt.Errorf("invalid congestion_level: %s", data.CongestionLevel)
		}
	}

	return nil
}