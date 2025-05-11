package models

import (
	"time"

	"github.com/google/uuid"
)

// CredentialAccessCountCalculator calculates the number of times credentials have been accessed
type CredentialAccessCountCalculator struct{}

// NewCredentialAccessCountCalculator creates a new instance of CredentialAccessCountCalculator
func NewCredentialAccessCountCalculator() *CredentialAccessCountCalculator {
	return &CredentialAccessCountCalculator{}
}

// Calculate computes the metric based on events
func (c *CredentialAccessCountCalculator) Calculate(events []Event) (Metric, error) {
	// Map to track access count by credential ID
	accessCounts := make(map[string]int)

	// Process events to count credential accesses
	for _, event := range events {
		// Handle password access events
		if event.Type == string(PasswordAccessedEvent) {
			// Extract passwordId from event data
			if data, ok := event.Data.(map[string]interface{}); ok {
				if id, exists := data["passwordId"]; exists {
					passwordId := id.(string)
					// Increment the access count for this credential
					accessCounts[passwordId]++
				}
			}
		}
	}

	// Calculate total access count
	totalAccessCount := 0
	for _, count := range accessCounts {
		totalAccessCount += count
	}

	return Metric{
		ID:          uuid.New().String(),
		Name:        "credential_access_count",
		Description: "Number of times credentials have been accessed",
		Value:       totalAccessCount,
		Unit:        "count",
		Timestamp:   time.Now(),
		Tags:        []string{string(PasswordMetric)},
	}, nil
}

// GetMetricType returns the type of metric
func (c *CredentialAccessCountCalculator) GetMetricType() MetricType {
	return CountMetric
}

// GetMetricName returns the name of the metric
func (c *CredentialAccessCountCalculator) GetMetricName() string {
	return "credential_access_count"
}
