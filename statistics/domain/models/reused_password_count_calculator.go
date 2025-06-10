package models

import (
	"time"

	"github.com/google/uuid"
)
// ReusedPasswordCountCalculator calculates the number of reused passwords
type ReusedPasswordCountCalculator struct{}

// NewReusedPasswordCountCalculator creates a new instance of ReusedPasswordCountCalculator
func NewReusedPasswordCountCalculator() *ReusedPasswordCountCalculator {
	return &ReusedPasswordCountCalculator{}
}

// Calculate computes the metric based on events
func (c *ReusedPasswordCountCalculator) Calculate(events []Event) (Metric, error) {
	// TODO
}

// GetMetricType returns the type of metric
func (c *ReusedPasswordCountCalculator) GetMetricType() MetricType {
	return CountMetric
}
// GetMetricName returns the name of the metric
func (c *ReusedPasswordCountCalculator) GetMetricName() string {
	return "reused_password_count"
}
// GetMetricDescription returns the description of the metric
func (c *ReusedPasswordCountCalculator) GetMetricDescription() string {
	return "Number of reused passwords"
}