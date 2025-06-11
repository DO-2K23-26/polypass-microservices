package models

import (
	"time"

	"github.com/google/uuid"
)

// OldPassword vérifie si un mot de passe est ancien (plus d'un an)
type OldPassword struct{}

// NewOldPassword crée une nouvelle instance de OldPassword
func NewOldPassword() *OldPassword {
	return &OldPassword{}
}

// Calculate calcule la métrique basée sur les événements
func (c *OldPassword) Calculate(events []Event) (Metric, error) {
	// Map pour stocker les mots de passe anciens et leurs IDs
	oldPasswords := make(map[string][]string)

	// Date limite (1 an avant aujourd'hui)
	oneYearAgo := time.Now().AddDate(-1, 0, 0)

	// Traitement des événements pour vérifier l'âge des mots de passe
	for _, event := range events {
		if event.Type == string(PasswordCreatedEvent) || event.Type == string(PasswordUpdatedEvent) {
			if data, ok := event.Data.(map[string]interface{}); ok {
				if password, exists := data["password"]; exists {
					if passwordStr, ok := password.(string); ok {
						if id, exists := data["passwordId"]; exists {
							if idStr, ok := id.(string); ok {
								if lastUpdated, exists := data["lastUpdated"]; exists {
									if lastUpdatedStr, ok := lastUpdated.(string); ok {
										lastUpdatedTime, err := time.Parse(time.RFC3339, lastUpdatedStr)
										if err == nil && lastUpdatedTime.Before(oneYearAgo) {
											oldPasswords[passwordStr] = append(oldPasswords[passwordStr], idStr)
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return Metric{
		ID:          uuid.New().String(),
		Name:        "old_password",
		Description: "Liste des mots de passe anciens (plus d'un an)",
		Value:       oldPasswords,
		Unit:        "list",
		Timestamp:   time.Now(),
		Tags:        []string{string(PasswordMetric)},
	}, nil
}

// GetMetricType retourne le type de métrique
func (c *OldPassword) GetMetricType() MetricType {
	return CountMetric
}

// GetMetricName retourne le nom de la métrique
func (c *OldPassword) GetMetricName() string {
	return "old_password"
}

// GetMetricDescription retourne la description de la métrique
func (c *OldPassword) GetMetricDescription() string {
	return "Liste des mots de passe anciens (plus d'un an)"
} 