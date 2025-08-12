package kafka

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

// Producer wraps a kafka.Writer.
type Producer struct {
	writer *kafka.Writer
}

// CreateProducer creates a new Kafka producer.
func CreateProducer(brokers []string, topic string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			WriteTimeout: 10 * time.Second,
			ReadTimeout:  10 * time.Second,
			ErrorLogger:  kafka.LoggerFunc(log.Printf),
		},
	}
}

// WriteMessages writes messages to Kafka.
func (p *Producer) WriteMessages(ctx context.Context, msgs ...kafka.Message) error {
	return p.writer.WriteMessages(ctx, msgs...)
}

// Close closes the Kafka writer.
func (p *Producer) Close() error {
	return p.writer.Close()
}
