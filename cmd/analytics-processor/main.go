package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"api-traffic-analytics/cmd/analytics-processor/internal/config"
	"api-traffic-analytics/cmd/analytics-processor/internal/repository"
	"api-traffic-analytics/cmd/analytics-processor/internal/service"
	"api-traffic-analytics/internal/interfaces"
	"api-traffic-analytics/internal/pkg/kafka"
	"api-traffic-analytics/internal/pkg/postgres"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize dependencies
	deps, err := initializeDependencies(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize dependencies: %v", err)
	}
	defer deps.Cleanup()

	// Create and start analytics service
	analyticsService := service.NewAnalyticsService(
		deps.KafkaConsumer,
		deps.AnalyticsProcessor,
		cfg,
	)

	// Handle graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	setupSignalHandler(cancel)

	// Start processing
	log.Println("Analytics Processor starting...")
	if err := analyticsService.Start(ctx); err != nil {
		log.Fatalf("Analytics service failed: %v", err)
	}
}

type Dependencies struct {
	KafkaConsumer      *kafka.Consumer
	AnalyticsProcessor interfaces.AnalyticsProcessor
}

func (d *Dependencies) Cleanup() {
	if d.KafkaConsumer != nil {
		d.KafkaConsumer.Close()
	}
	// Cleanup other resources
}

func initializeDependencies(cfg *config.Config) (*Dependencies, error) {
	// Initialize Kafka Consumer
	consumer := kafka.CreateConsumer(cfg.KafkaBrokers, cfg.KafkaTopic, cfg.KafkaGroupID)

	// Initialize Database
	db, err := postgres.ConnectDB()
	if err != nil {
		return nil, err
	}

	// Initialize repositories
	postgresRepo := postgres.NewTrafficDataRepository(db)
	repo := repository.NewRepository(postgresRepo)

	// Initialize processors
	analyticsProcessor := service.NewAnalyticsProcessor(repo)

	return &Dependencies{
		KafkaConsumer:      consumer,
		AnalyticsProcessor: analyticsProcessor,
	}, nil
}

func setupSignalHandler(cancel context.CancelFunc) {
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigterm
		log.Println("Shutdown signal received...")
		cancel()
	}()
}
