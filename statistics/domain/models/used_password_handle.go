package models

import (
	"time"

	"github.com/google/uuid"
)

// PasswordUsageCountCalculator counts the number of times a password has been used
type PasswordUsageCountCalculator struct{}

// NewPasswordUsageCountCalculator creates a new instance
func NewPasswordUsageCountCalculator() *PasswordUsageCountCalculator {
	return &PasswordUsageCountCalculator{}
}

// Calculate computes the number of times a specific password has been used
func (c *PasswordUsageCountCalculator) Calculate(events []Event, passwordID string) (Metric, error) {
	count := 0

	for _, event := range events {
		if event.Type == string(PasswordUsedEvent) {
			if data, ok := event.Data.(map[string]interface{}); ok {
				if id, exists := data["passwordId"]; exists && id == passwordID {
					count++
				}
			}
		}
	}

	return Metric{
		ID:          uuid.New().String(),
		Name:        "password_usage_count",
		Description: "Number of times the password has been used",
		Value:       count,
		Unit:        "count",
		Timestamp:   time.Now(),
		Tags:        []string{string(PasswordMetric)},
	}, nil
}

// GetMetricType returns the type of metric
func (c *PasswordUsageCountCalculator) GetMetricType() MetricType {
	return CountMetric
}

// GetMetricName returns the name of the metric
func (c *PasswordUsageCountCalculator) GetMetricName() string {
	return "password_usage_count"
}
