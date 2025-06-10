package models

import (
	"time"

	"github.com/google/uuid"
)

//WeakPasswordCountCalculator calculates the number of weak passwords
type WeakPasswordCountCalculator struct{}

// NewWeakPasswordCountCalculator creates a new instance of WeakPasswordCountCalculator
func NewWeakPasswordCountCalculator() *WeakPasswordCountCalculator {
	return &WeakPasswordCountCalculator{}
}

// Calculate computes the metric based on events
func (c *WeakPasswordCountCalculator) Calculate(events []Event) (Metric, error) {
	//TODO
}

// GetMetricType returns the type of metric
func (c *WeakPasswordCountCalculator) GetMetricType() MetricType {
	return CountMetric
}

// GetMetricName returns the name of the metric
func (c *WeakPasswordCountCalculator) GetMetricName() string {
	return "weak_password_count"
}
// GetMetricDescription returns the description of the metric
func (c *WeakPasswordCountCalculator) GetMetricDescription() string {
	return "Number of weak passwords"
}