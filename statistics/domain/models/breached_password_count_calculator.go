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
	// Map to track breached passwords by ID
	breachedPasswords := make(map[string]int)

	// Process events to count breached passwords
	for _, event := range events {
		// Handle password breach events
		if event.Type == string(PasswordBreachEvent) {
			// Extract passwordId from event data
			if data, ok := event.Data.(map[string]interface{}); ok {
				if id, exists := data["passwordId"]; exists {
					passwordId := id.(string)
					// Mark this password as breached
					// TODO: Implement logic to check if the password is actually breached
					breachedPasswords[passwordId]++
				}
			}
		}
	}

	// Count the number of breached passwords
	count := len(breachedPasswords)

	return Metric{
		ID:          uuid.New().String(),
		Name:        "breached_password_count",
		Description: "Number of breached passwords",
		Value:       count,
		Unit:        "count",
		Timestamp:   time.Now(),
		Tags:        []string{string(PasswordMetric)},
	}, nil
}

// GetMetricType returns the type of metric
func (c *BreachedPasswordCountCalculator) GetMetricType() MetricType {
	return CountMetric
}

// GetMetricName returns the name of the metric
func (c *BreachedPasswordCountCalculator) GetMetricName() string {
	return "breached_password_count"
}
