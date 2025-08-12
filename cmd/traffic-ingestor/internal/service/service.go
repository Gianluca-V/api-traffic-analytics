package service

import (
	"context"
	"encoding/json"
	"log"

	"api-traffic-analytics/cmd/traffic-ingestor/internal/repository"
	kafkaPkg "api-traffic-analytics/internal/pkg/kafka"
	"api-traffic-analytics/internal/shared/models"
	"github.com/segmentio/kafka-go"
)

type Service struct {
	repo     *repository.Repository
	producer *kafkaPkg.Producer
}

func NewService(repo *repository.Repository, producer *kafkaPkg.Producer) *Service {
	return &Service{repo: repo, producer: producer}
}

func (s *Service) ProcessTrafficData(ctx context.Context, data *models.TrafficData) error {
	// Store data
	if err := s.repo.StoreTrafficData(ctx, data); err != nil {
		return err
	}

	// Publish to Kafka
	msg, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = s.producer.WriteMessages(ctx, kafka.Message{
		Value: msg,
	})
	if err != nil {
		log.Printf("failed to write message to kafka: %v", err)
		return err
	}

	log.Println("Successfully published message to kafka")
	return nil
}