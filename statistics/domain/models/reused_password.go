package models

import (
	"time"

	"github.com/google/uuid"
)

// ReusedPassword calculates the number of reused passwords
type ReusedPassword struct{}

// NewReusedPassword creates a new instance of ReusedPassword
func NewReusedPassword() *ReusedPassword {
	return &ReusedPassword{}
}

// Calculate computes the metric based on events
func (c *ReusedPassword) Calculate(events []Event) (Metric, error) {
	// Map pour stocker les mots de passe et leurs occurrences
	passwordOccurrences := make(map[string][]string)

	// Traitement des événements pour détecter les mots de passe réutilisés
	for _, event := range events {
		if event.Type == string(PasswordCreatedEvent) || event.Type == string(PasswordUpdatedEvent) {
			if data, ok := event.Data.(map[string]interface{}); ok {
				if password, exists := data["password"]; exists {
					if passwordStr, ok := password.(string); ok {
						if id, exists := data["passwordId"]; exists {
							if idStr, ok := id.(string); ok {
								passwordOccurrences[passwordStr] = append(passwordOccurrences[passwordStr], idStr)
							}
						}
					}
				}
			}
		}
	}

	// Filtrer pour ne garder que les mots de passe réutilisés
	reusedPasswords := make(map[string][]string)
	for password, ids := range passwordOccurrences {
		if len(ids) > 1 {
			reusedPasswords[password] = ids
		}
	}

	return Metric{
		ID:          uuid.New().String(),
		Name:        "reused_password",
		Description: "List of reused passwords",
		Value:       reusedPasswords,
		Unit:        "list",
		Timestamp:   time.Now(),
		Tags:        []string{string(PasswordMetric)},
	}, nil
}

// GetMetricType returns the type of metric
func (c *ReusedPassword) GetMetricType() MetricType {
	return CountMetric
}
// GetMetricName returns the name of the metric
func (c *ReusedPassword) GetMetricName() string {
	return "reused_password"
}
// GetMetricDescription returns the description of the metric
func (c *ReusedPassword) GetMetricDescription() string {
	return "List of reused passwords"
}