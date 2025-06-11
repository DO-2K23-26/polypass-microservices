package models

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

// BreachedPassword vérifie si un mot de passe a été compromis
type BreachedPassword struct{}

// NewBreachedPassword crée une nouvelle instance de BreachedPassword
func NewBreachedPassword() *BreachedPassword {
	return &BreachedPassword{}
}

// Calculate calcule la métrique basée sur les événements
func (c *BreachedPassword) Calculate(events []Event) (Metric, error) {
	// Map pour stocker les mots de passe compromis et leurs IDs
	breachedPasswords := make(map[string][]string)

	// Traitement des événements pour vérifier les mots de passe compromis
	for _, event := range events {
		if event.Type == string(PasswordCreatedEvent) || event.Type == string(PasswordUpdatedEvent) {
			if data, ok := event.Data.(map[string]interface{}); ok {
				if password, exists := data["password"]; exists {
					if passwordStr, ok := password.(string); ok {
						if id, exists := data["passwordId"]; exists {
							if idStr, ok := id.(string); ok {
								// Vérifier si le mot de passe est compromis
								if isBreached, err := CheckPasswordBreach(passwordStr); err == nil && isBreached {
									breachedPasswords[passwordStr] = append(breachedPasswords[passwordStr], idStr)
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
		Name:        "breached_password",
		Description: "Liste des mots de passe compromis",
		Value:       breachedPasswords,
		Unit:        "list",
		Timestamp:   time.Now(),
		Tags:        []string{string(PasswordMetric)},
	}, nil
}

// GetMetricType retourne le type de métrique
func (c *BreachedPassword) GetMetricType() MetricType {
	return CountMetric
}

// GetMetricName retourne le nom de la métrique
func (c *BreachedPassword) GetMetricName() string {
	return "breached_password"
}

// GetMetricDescription retourne la description de la métrique
func (c *BreachedPassword) GetMetricDescription() string {
	return "Liste des mots de passe compromis"
}

// CheckPasswordBreach vérifie si un mot de passe a été compromis en utilisant l'API HIBP
func CheckPasswordBreach(password string) (bool, error) {
	// Calculer le hash SHA-1 du mot de passe
	hash := sha1.New()
	io.WriteString(hash, password)
	hashBytes := hash.Sum(nil)
	hashString := strings.ToUpper(hex.EncodeToString(hashBytes))

	// Extraire les 5 premiers caractères du hash
	prefix := hashString[:5]
	suffix := hashString[5:]

	// Appeler l'API HIBP
	url := fmt.Sprintf("https://api.pwnedpasswords.com/range/%s", prefix)
	resp, err := http.Get(url)
	if err != nil {
		return false, fmt.Errorf("erreur lors de l'appel à l'API HIBP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("erreur de l'API HIBP: %d", resp.StatusCode)
	}

	// Lire la réponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("erreur lors de la lecture de la réponse: %w", err)
	}

	// Vérifier si le suffixe du hash est présent dans la réponse
	lines := strings.Split(string(body), "\r\n")
	for _, line := range lines {
		if strings.HasPrefix(line, suffix) {
			return true, nil
		}
	}

	return false, nil
} 