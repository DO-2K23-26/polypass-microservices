package models

import (
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
)

// StrongPassword vérifie si un mot de passe est fort
type StrongPassword struct{}

// NewStrongPassword crée une nouvelle instance de StrongPassword
func NewStrongPassword() *StrongPassword {
	return &StrongPassword{}
}

// Calculate calcule la métrique basée sur les événements
func (c *StrongPassword) Calculate(events []Event) (Metric, error) {
	// Map pour stocker les mots de passe forts et leurs IDs
	strongPasswords := make(map[string][]string)

	// Traitement des événements pour vérifier la force des mots de passe
	for _, event := range events {
		if event.Type == string(PasswordCreatedEvent) || event.Type == string(PasswordUpdatedEvent) {
			if data, ok := event.Data.(map[string]interface{}); ok {
				if password, exists := data["password"]; exists {
					if passwordStr, ok := password.(string); ok {
						if id, exists := data["passwordId"]; exists {
							if idStr, ok := id.(string); ok {
								// Si le mot de passe est fort, l'ajouter à la map
								if IsStrongPassword(passwordStr) {
									strongPasswords[passwordStr] = append(strongPasswords[passwordStr], idStr)
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
		Name:        "strong_password",
		Description: "Liste des mots de passe forts",
		Value:       strongPasswords,
		Unit:        "list",
		Timestamp:   time.Now(),
		Tags:        []string{string(PasswordMetric)},
	}, nil
}

// GetMetricType retourne le type de métrique
func (c *StrongPassword) GetMetricType() MetricType {
	return CountMetric
}

// GetMetricName retourne le nom de la métrique
func (c *StrongPassword) GetMetricName() string {
	return "strong_password"
}

// GetMetricDescription retourne la description de la métrique
func (c *StrongPassword) GetMetricDescription() string {
	return "Liste des mots de passe forts"
}

// IsStrongPassword vérifie si un mot de passe est fort
func IsStrongPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case strings.ContainsRune("@$!%*?&", char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}
