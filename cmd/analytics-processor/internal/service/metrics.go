package service

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	messagesProcessed prometheus.Counter
	messagesFailed    prometheus.Counter
	processingTime    prometheus.Histogram
}

func NewMetrics() *Metrics {
	return &Metrics{
		messagesProcessed: promauto.NewCounter(prometheus.CounterOpts{
			Name: "analytics_processor_messages_processed_total",
			Help: "Total messages processed",
		}),
		messagesFailed: promauto.NewCounter(prometheus.CounterOpts{
			Name: "analytics_processor_messages_failed_total",
			Help: "Total messages failed to process",
		}),
		processingTime: promauto.NewHistogram(prometheus.HistogramOpts{
			Name: "analytics_processor_processing_duration_seconds",
			Help: "Processing duration in seconds",
		}),
	}
}

func (m *Metrics) IncrementProcessed() {
	m.messagesProcessed.Inc()
}

func (m *Metrics) IncrementFailed() {
	m.messagesFailed.Inc()
}

func (m *Metrics) RecordProcessingTime(duration time.Duration) {
	m.processingTime.Observe(duration.Seconds())
}
