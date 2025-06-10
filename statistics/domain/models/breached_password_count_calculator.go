package models

import (
	"time"

	"github.com/google/uuid"
)

// BreachedPasswordCountCalculator calculates the number of breached passwords
type BreachedPasswordCountCalculator struct{}

// NewBreachedPasswordCountCalculator creates a new instance of BreachedPasswordCountCalculator
func NewBreachedPasswordCountCalculator() *BreachedPasswordCountCalculator {
	return &BreachedPasswordCountCalculator{}
}

// Calculate computes the metric based on events
func (c *BreachedPasswordCountCalculator) Calculate(events []Event) (Metric, error) {
// TODO
}

// GetMetricType returns the type of metric
func (c *BreachedPasswordCountCalculator) GetMetricType() MetricType {
	return CountMetric
}

// GetMetricName returns the name of the metric
func (c *BreachedPasswordCountCalculator) GetMetricName() string {
	return "breached_password_count"
}
// GetMetricDescription returns the description of the metric
func (c *BreachedPasswordCountCalculator) GetMetricDescription() string {
	return "Number of breached passwords"
}
