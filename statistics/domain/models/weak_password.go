package models

import (
	"time"

	"github.com/google/uuid"
)

// WeakPassword vérifie si un mot de passe est faible
type WeakPassword struct{}

// NewWeakPassword crée une nouvelle instance de WeakPassword
func NewWeakPassword() *WeakPassword {
	return &WeakPassword{}
}

// Calculate calcule la métrique basée sur les événements
func (c *WeakPassword) Calculate(events []Event) (Metric, error) {
	// Map pour stocker les mots de passe faibles et leurs IDs
	weakPasswords := make(map[string][]string)

	// Traitement des événements pour vérifier la force des mots de passe
	for _, event := range events {
		if event.Type == string(PasswordCreatedEvent) || event.Type == string(PasswordUpdatedEvent) {
			if data, ok := event.Data.(map[string]interface{}); ok {
				if password, exists := data["password"]; exists {
					if passwordStr, ok := password.(string); ok {
						if id, exists := data["passwordId"]; exists {
							if idStr, ok := id.(string); ok {
								// Si le mot de passe n'est pas fort, il est considéré comme faible
								if !IsStrongPassword(passwordStr) {
									weakPasswords[passwordStr] = append(weakPasswords[passwordStr], idStr)
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
		Name:        "weak_password",
		Description: "Liste des mots de passe faibles",
		Value:       weakPasswords,
		Unit:        "list",
		Timestamp:   time.Now(),
		Tags:        []string{string(PasswordMetric)},
	}, nil
}

// GetMetricType retourne le type de métrique
func (c *WeakPassword) GetMetricType() MetricType {
	return CountMetric
}

// GetMetricName retourne le nom de la métrique
func (c *WeakPassword) GetMetricName() string {
	return "weak_password"
}

// GetMetricDescription retourne la description de la métrique
func (c *WeakPassword) GetMetricDescription() string {
	return "Liste des mots de passe faibles"
}

