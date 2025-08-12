package service

import (
	"context"
	"log"
	"time"

	"api-traffic-analytics/cmd/analytics-processor/internal/config"
	handler "api-traffic-analytics/cmd/analytics-processor/internal/handler"
	"api-traffic-analytics/internal/interfaces"
	"api-traffic-analytics/internal/pkg/kafka"
)

type AnalyticsService struct {
	consumer       *kafka.Consumer
	processor      interfaces.AnalyticsProcessor // Cambiado a la interfaz
	messageHandler *handler.MessageHandler
	config         *config.Config
	metrics        *Metrics
	logger         *log.Logger
}

func NewAnalyticsService(
	consumer *kafka.Consumer,
	processor interfaces.AnalyticsProcessor, // Cambiado a la interfaz
	config *config.Config,
) *AnalyticsService {
	logger := log.Default()

	return &AnalyticsService{
		consumer:       consumer,
		processor:      processor,
		messageHandler: handler.NewMessageHandler(processor, logger),
		config:         config,
		metrics:        NewMetrics(),
		logger:         logger,
	}
}

func (s *AnalyticsService) Start(ctx context.Context) error {
	s.logger.Println("Starting analytics processing loop...")

	for {
		select {
		case <-ctx.Done():
			s.logger.Println("Analytics service stopping...")
			return nil
		default:
			if err := s.processNextMessage(ctx); err != nil {
				s.logger.Printf("Error processing message: %v", err)
				time.Sleep(1 * time.Second)
			}
		}
	}
}

func (s *AnalyticsService) processNextMessage(ctx context.Context) error {
	// Add timeout to prevent hanging
	msgCtx, cancel := context.WithTimeout(ctx, time.Duration(s.config.ProcessingTimeout)*time.Second)
	defer cancel()

	msg, err := s.consumer.FetchMessage(msgCtx)
	if err != nil {
		return err
	}

	// Process message using handler
	startTime := time.Now()
	err = s.messageHandler.HandleMessage(msgCtx, msg.Value)

	// Create and log metadata
	processingTime := time.Since(startTime)
	// metadata := s.messageHandler.CreateMessageMetadata(msg, processingTime, err == nil, err)
	// s.messageHandler.LogMessageMetadata(metadata)

	if err != nil {
		s.metrics.IncrementFailed()
		// Handle error (opcional)
		// s.messageHandler.HandleError(msgCtx, msg, err)
		return err
	}

	s.metrics.IncrementProcessed()
	s.metrics.RecordProcessingTime(processingTime)

	// Commit message only after successful processing
	if err := s.consumer.CommitMessages(ctx, msg); err != nil {
		s.logger.Printf("Warning: Error committing message: %v", err)
	}

	return nil
}