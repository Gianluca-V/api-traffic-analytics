package config

import (
	"os"
	"strconv"
)

type Config struct {
	KafkaBrokers      []string
	KafkaTopic        string
	KafkaGroupID      string
	MetricsPort       string
	ProcessingTimeout int
	BatchSize         int
}

func Load() *Config {
	return &Config{
		KafkaBrokers:      []string{getEnv("KAFKA_BROKER", "kafka:9092")},
		KafkaTopic:        getEnv("KAFKA_TOPIC_TRAFFIC", "traffic-data"),
		KafkaGroupID:      getEnv("KAFKA_CONSUMER_GROUP", "analytics-processor"),
		MetricsPort:       getEnv("METRICS_PORT", "8080"),
		ProcessingTimeout: getIntEnv("PROCESSING_TIMEOUT", 30),
		BatchSize:         getIntEnv("BATCH_SIZE", 1),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		res, _ := strconv.Atoi(value)
		return res
	}
	return defaultValue
}
