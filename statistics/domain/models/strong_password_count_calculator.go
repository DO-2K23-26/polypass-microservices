package models

import (
	"time"

	"github.com/google/uuid"
)

// StrongPasswordCountCalculator calculates the number of strong passwords
type StrongPasswordCountCalculator struct{}
// NewStrongPasswordCountCalculator creates a new instance of StrongPasswordCountCalculator
func NewStrongPasswordCountCalculator() *StrongPasswordCountCalculator {
	return &StrongPasswordCountCalculator{}
}

// Calculate computes the metric based on events
func (c *StrongPasswordCountCalculator) Calculate(events []Event) (Metric, error) {
	//TODO
}

// GetMetricType returns the type of metric
func (c *StrongPasswordCountCalculator) GetMetricType() MetricType {
	return CountMetric
}
// GetMetricName returns the name of the metric
func (c *StrongPasswordCountCalculator) GetMetricName() string {
	return "strong_password_count"
}
// GetMetricDescription returns the description of the metric
func (c *StrongPasswordCountCalculator) GetMetricDescription() string {
	return "Number of strong passwords"
}