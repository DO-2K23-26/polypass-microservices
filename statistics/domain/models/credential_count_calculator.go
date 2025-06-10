package models

import (
	"time"

	"github.com/google/uuid"
)

// CredentialCountCalculator calculates the number of credentials a user has
type CredentialCountCalculator struct{}

// NewCredentialCountCalculator creates a new instance of CredentialCountCalculator
func NewCredentialCountCalculator() *CredentialCountCalculator {
	return &CredentialCountCalculator{}
}

// Calculate computes the metric based on events
func (c *CredentialCountCalculator) Calculate(events []Event) (Metric, error) {
	// Map to track credentials by ID
	credentials := make(map[string]bool)

	// Process events to count credentials
	for _, event := range events {
		// Extract passwordId from event data
		var passwordId string

		// Handle different event types
		switch event.Type {
		case string(PasswordCreatedEvent):
			// When a password is created, add it to the count
			if data, ok := event.Data.(map[string]interface{}); ok {
				if id, exists := data["passwordId"]; exists {
					passwordId = id.(string)
					credentials[passwordId] = true
				}
			}
		case string(PasswordDeletedEvent):
			// When a password is deleted, remove it from the count
			if data, ok := event.Data.(map[string]interface{}); ok {
				if id, exists := data["passwordId"]; exists {
					passwordId = id.(string)
					delete(credentials, passwordId)
				}
			}
		}
	}

	// Count the number of credentials
	count := len(credentials)

	return Metric{
		ID:          uuid.New().String(),
		Name:        "credential_count",
		Description: "Number of credentials",
		Value:       count,
		Unit:        "count",
		Timestamp:   time.Now(),
		Tags:        []string{string(PasswordMetric)},
	}, nil
}

// GetMetricType returns the type of metric
func (c *CredentialCountCalculator) GetMetricType() MetricType {
	return CountMetric
}

// GetMetricName returns the name of the metric
func (c *CredentialCountCalculator) GetMetricName() string {
	return "credential_count"
}
