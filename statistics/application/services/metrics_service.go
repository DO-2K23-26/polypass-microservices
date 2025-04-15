package services

import (
	"context"
	"encoding/json"
	"time"

	"github.com/DO-2K23-26/polypass-microservices/statistics/domain/models"
	"github.com/DO-2K23-26/polypass-microservices/statistics/domain/repositories"
)

// MetricsService handles all metrics-related operations
type MetricsService struct {
	repo repositories.EventSourcingRepository
}

// NewMetricsService creates a new MetricsService
func NewMetricsService(repo repositories.EventSourcingRepository) *MetricsService {
	return &MetricsService{
		repo: repo,
	}
}

// ProcessCredentialEvent processes credential-related events
func (s *MetricsService) ProcessCredentialEvent(ctx context.Context, data []byte) error {
	var eventData struct {
		UserID       string `json:"user_id"`
		GroupID      string `json:"group_id"`
		CredentialID string `json:"credential_id"`
		Action       string `json:"action"`
		Timestamp    string `json:"timestamp"`
	}

	if err := json.Unmarshal(data, &eventData); err != nil {
		return err
	}

	// Store the event
	if err := s.repo.StoreEvent(ctx, repositories.Event{
		Type:      repositories.CredentialEvent,
		Data:      data,
		Timestamp: time.Now(),
	}); err != nil {
		return err
	}

	return nil
}

// ProcessAccessEvent processes credential access events
func (s *MetricsService) ProcessAccessEvent(ctx context.Context, data []byte) error {
	var eventData struct {
		CredentialID string `json:"credential_id"`
		UserID       string `json:"user_id"`
		GroupID      string `json:"group_id"`
		IPAddress    string `json:"ip_address"`
		UserAgent    string `json:"user_agent"`
		IsOneTime    bool   `json:"is_one_time"`
		SharedAt     string `json:"shared_at"`
	}

	if err := json.Unmarshal(data, &eventData); err != nil {
		return err
	}

	// Store the event
	if err := s.repo.StoreEvent(ctx, repositories.Event{
		Type:      repositories.AccessEvent,
		Data:      data,
		Timestamp: time.Now(),
	}); err != nil {
		return err
	}

	return nil
}

// GetUserMetrics retrieves user metrics by replaying events
func (s *MetricsService) GetUserMetrics(ctx context.Context, userID string) (*models.UserMetrics, error) {
	return s.repo.GetUserMetrics(ctx, userID)
}

// GetGroupMetrics retrieves group metrics by replaying events
func (s *MetricsService) GetGroupMetrics(ctx context.Context, groupID string) (*models.GroupMetrics, error) {
	return s.repo.GetGroupMetrics(ctx, groupID)
}

// GetCredentialAccesses retrieves credential access records by replaying events
func (s *MetricsService) GetCredentialAccesses(ctx context.Context, credentialID string, startDate, endDate time.Time) ([]models.CredentialAccess, error) {
	return s.repo.GetCredentialAccesses(ctx, credentialID, startDate, endDate)
}

// GetSharedCredentialStats retrieves shared credential statistics
func (s *MetricsService) GetSharedCredentialStats(ctx context.Context, credentialID string) (*models.SharedCredentialStats, error) {
	return s.repo.GetSharedCredentialStats(ctx, credentialID)
}

// GetCredentialTrends retrieves credential creation and usage trends
func (s *MetricsService) GetCredentialTrends(ctx context.Context, startDate, endDate time.Time) (*models.CredentialTrend, error) {
	return s.repo.GetCredentialTrends(ctx, startDate, endDate)
}

// GetPasswordStrengths retrieves password strength analysis
func (s *MetricsService) GetPasswordStrengths(ctx context.Context, userID string) ([]models.PasswordStrength, error) {
	return s.repo.GetPasswordStrengths(ctx, userID)
}

// GetReusedPasswords retrieves reused password analysis
func (s *MetricsService) GetReusedPasswords(ctx context.Context, userID string) ([]models.ReusedPassword, error) {
	return s.repo.GetReusedPasswords(ctx, userID)
}

// GetBreachedCredentials retrieves breached credential analysis
func (s *MetricsService) GetBreachedCredentials(ctx context.Context, userID string) ([]models.BreachedCredential, error) {
	return s.repo.GetBreachedCredentials(ctx, userID)
}

// GetOldPasswords retrieves old password analysis
func (s *MetricsService) GetOldPasswords(ctx context.Context, userID string, daysThreshold int) ([]models.OldPassword, error) {
	return s.repo.GetOldPasswords(ctx, userID, daysThreshold)
}
