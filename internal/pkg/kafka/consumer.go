package kafka

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

// Consumer wraps a kafka.Reader.
type Consumer struct {
	reader *kafka.Reader
}

// CreateConsumer creates a new Kafka consumer.
func CreateConsumer(brokers []string, topic string, groupID string) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:        brokers,
			GroupID:        groupID,
			Topic:          topic,
			MinBytes:       10e3, // 10KB
			MaxBytes:       10e6, // 10MB
			CommitInterval: time.Second,
			ErrorLogger:    kafka.LoggerFunc(log.Printf),
		}),
	}
}

// FetchMessage fetches a message from Kafka.
func (c *Consumer) FetchMessage(ctx context.Context) (kafka.Message, error) {
	return c.reader.FetchMessage(ctx)
}

// CommitMessages commits messages.
func (c *Consumer) CommitMessages(ctx context.Context, msgs ...kafka.Message) error {
	return c.reader.CommitMessages(ctx, msgs...)
}

// Close closes the Kafka reader.
func (c *Consumer) Close() error {
	return c.reader.Close()
}
