package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port                string
	Environment         string
	APIKey              string
	RateLimitRequests   int
	RateLimitDuration   int
	TrafficIngestorURL  string
	AnalyticsServiceURL string
	AlertingServiceURL  string
}

func Load() *Config {
	rateLimitRequests, _ := strconv.Atoi(getEnv("RATE_LIMIT_REQUESTS", "100"))
	rateLimitDuration, _ := strconv.Atoi(getEnv("RATE_LIMIT_DURATION", "60"))

	return &Config{
		Port:                getEnv("PORT", "8080"),
		Environment:         getEnv("ENVIRONMENT", "development"),
		APIKey:              getEnv("API_KEY", "default-api-key-change-in-production"),
		RateLimitRequests:   rateLimitRequests,
		RateLimitDuration:   rateLimitDuration,
		TrafficIngestorURL:  getEnv("TRAFFIC_INGESTOR_URL", "http://traffic-ingestor:8081"),
		AnalyticsServiceURL: getEnv("ANALYTICS_SERVICE_URL", "http://analytics-processor:8082"),
		AlertingServiceURL:  getEnv("ALERTING_SERVICE_URL", "http://alerting-service:8083"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}