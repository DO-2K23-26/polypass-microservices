package models

import (
	"time"
)

// Metric represents a statistical metric derived from events
type Metric struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Value       interface{} `json:"value"`
	Unit        string      `json:"unit"`
	Timestamp   time.Time   `json:"timestamp"`
	Tags        []string    `json:"tags"`
}

// MetricType defines the possible types of metrics
type MetricType string

const (
	// Define common metric types
	CountMetric     MetricType = "COUNT"
	GaugeMetric     MetricType = "GAUGE"
	HistogramMetric MetricType = "HISTOGRAM"
	SummaryMetric   MetricType = "SUMMARY"
	// Add more metric types as needed
)

// MetricCategory defines the categories of metrics
type MetricCategory string

const (
	// Define common metric categories
	UserMetric        MetricCategory = "USER"
	PasswordMetric    MetricCategory = "PASSWORD"
	PerformanceMetric MetricCategory = "PERFORMANCE"
	SecurityMetric    MetricCategory = "SECURITY"
	// Add more categories as needed
)

// MetricCalculator defines the interface for calculating metrics
type MetricCalculator interface {
	// Calculate computes a metric based on a collection of events
	Calculate(events []Event) (Metric, error)

	// GetMetricType returns the type of metric this calculator produces
	GetMetricType() MetricType

	// GetMetricName returns the name of the metric
	GetMetricName() string
}
