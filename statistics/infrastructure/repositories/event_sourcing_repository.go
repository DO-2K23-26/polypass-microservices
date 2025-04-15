package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/DO-2K23-26/polypass-microservices/statistics/domain/models"
	"github.com/DO-2K23-26/polypass-microservices/statistics/domain/repositories"
)

const (
	CredentialEvent repositories.EventType = "credential"
	AccessEvent     repositories.EventType = "access"
)

// EventSourcingRepository implements the event sourcing pattern
type EventSourcingRepository struct {
	events []repositories.Event
	mu     sync.RWMutex
}

// NewEventSourcingRepository creates a new event sourcing repository
func NewEventSourcingRepository() *EventSourcingRepository {
	return &EventSourcingRepository{
		events: make([]repositories.Event, 0),
	}
}

// Setup initializes the repository
func (r *EventSourcingRepository) Setup() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.events = make([]repositories.Event, 0)
	return nil
}

// Shutdown cleans up resources
func (r *EventSourcingRepository) Shutdown() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.events = nil
	return nil
}

// StoreEvent stores a new event
func (r *EventSourcingRepository) StoreEvent(ctx context.Context, event repositories.Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.events = append(r.events, event)
	return nil
}

// GetEventsByType retrieves events of a specific type
func (r *EventSourcingRepository) GetEventsByType(ctx context.Context, eventType repositories.EventType) ([]repositories.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filtered []repositories.Event
	for _, event := range r.events {
		if event.Type == repositories.EventType(eventType) {
			filtered = append(filtered, event)
		}
	}

	return filtered, nil
}

// GetEventsByTimeRange retrieves events within a time range
func (r *EventSourcingRepository) GetEventsByTimeRange(ctx context.Context, start, end time.Time) ([]repositories.Event, error) {
	if start.After(end) {
		return nil, ErrInvalidTimeRange
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	var filtered []repositories.Event
	for _, event := range r.events {
		if event.Timestamp.After(start) && event.Timestamp.Before(end) {
			filtered = append(filtered, event)
		}
	}

	return filtered, nil
}

// GetUserMetrics calculates user metrics from events
func (r *EventSourcingRepository) GetUserMetrics(ctx context.Context, userID string) (*models.UserMetrics, error) {
	if userID == "" {
		return nil, ErrInvalidUserID
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	metrics := &models.UserMetrics{
		UserID:      userID,
		LastUpdated: time.Now(),
	}

	// Get all credential events for this user
	credentialEvents, err := r.GetEventsByType(ctx, repositories.EventType(CredentialEvent))
	if err != nil {
		return nil, err
	}

	for _, event := range credentialEvents {
		var data struct {
			UserID string `json:"user_id"`
			Action string `json:"action"`
		}

		if err := json.Unmarshal(event.Data, &data); err != nil {
			continue
		}

		if data.UserID != userID {
			continue
		}

		switch data.Action {
		case "create":
			metrics.TotalCredentials++
		case "delete":
			metrics.TotalCredentials--
		}
	}

	return metrics, nil
}

// GetGroupMetrics calculates group metrics from events
func (r *EventSourcingRepository) GetGroupMetrics(ctx context.Context, groupID string) (*models.GroupMetrics, error) {
	if groupID == "" {
		return nil, ErrInvalidGroupID
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	metrics := &models.GroupMetrics{
		GroupID:     groupID,
		LastUpdated: time.Now(),
	}

	// Get all credential events for this group
	credentialEvents, err := r.GetEventsByType(ctx, repositories.EventType(CredentialEvent))
	if err != nil {
		return nil, err
	}

	uniqueUsers := make(map[string]bool)

	for _, event := range credentialEvents {
		var data struct {
			GroupID string `json:"group_id"`
			UserID  string `json:"user_id"`
			Action  string `json:"action"`
		}

		if err := json.Unmarshal(event.Data, &data); err != nil {
			continue
		}

		if data.GroupID != groupID {
			continue
		}

		uniqueUsers[data.UserID] = true

		switch data.Action {
		case "create":
			metrics.TotalCredentials++
		case "delete":
			metrics.TotalCredentials--
		}
	}

	metrics.ActiveUsers = len(uniqueUsers)
	return metrics, nil
}

// GetSharedCredentialStats retrieves shared credential statistics
func (r *EventSourcingRepository) GetSharedCredentialStats(ctx context.Context, credentialID string) (*models.SharedCredentialStats, error) {
	if credentialID == "" {
		return nil, ErrInvalidCredentialID
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	// Get all access events for this credential
	accessEvents, err := r.GetEventsByType(ctx, repositories.EventType(AccessEvent))
	if err != nil {
		return nil, err
	}

	stats := &models.SharedCredentialStats{
		CredentialID: credentialID,
	}

	uniqueViewers := make(map[string]bool)
	oneTimeViews := make(map[string]bool)

	for _, event := range accessEvents {
		var accessData struct {
			CredentialID string `json:"credential_id"`
			UserID       string `json:"user_id"`
			IsOneTime    bool   `json:"is_one_time"`
			Timestamp    string `json:"timestamp"`
		}

		if err := json.Unmarshal(event.Data, &accessData); err != nil {
			continue
		}

		if accessData.CredentialID != credentialID {
			continue
		}

		stats.TotalViews++
		uniqueViewers[accessData.UserID] = true

		if accessData.IsOneTime {
			oneTimeViews[accessData.UserID] = true
		}

		timestamp, err := time.Parse(time.RFC3339, accessData.Timestamp)
		if err != nil {
			continue
		}

		if stats.FirstShared.IsZero() || timestamp.Before(stats.FirstShared) {
			stats.FirstShared = timestamp
		}

		if timestamp.After(stats.LastAccessed) {
			stats.LastAccessed = timestamp
		}
	}

	stats.UniqueViewers = len(uniqueViewers)
	stats.OneTimeViews = len(oneTimeViews)

	return stats, nil
}

// GetCredentialTrends retrieves credential creation and usage trends
func (r *EventSourcingRepository) GetCredentialTrends(ctx context.Context, startDate, endDate time.Time) (*models.CredentialTrend, error) {
	if startDate.After(endDate) {
		return nil, ErrInvalidTimeRange
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	// Get all events within the time range
	events, err := r.GetEventsByTimeRange(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}

	trend := &models.CredentialTrend{
		StartDate: startDate,
		EndDate:   endDate,
	}

	// Group events by day
	eventsByDay := make(map[string]struct {
		creations int
		accesses  int
	})

	for _, event := range events {
		day := event.Timestamp.Format("2006-01-02")

		dayStats := eventsByDay[day]
		if event.Type == repositories.EventType(CredentialEvent) {
			dayStats.creations++
		} else if event.Type == repositories.EventType(AccessEvent) {
			dayStats.accesses++
		}
		eventsByDay[day] = dayStats
	}

	// Convert to data points
	for day, stats := range eventsByDay {
		date, _ := time.Parse("2006-01-02", day)
		trend.DataPoints = append(trend.DataPoints, models.TrendDataPoint{
			Date:      date,
			Creations: stats.creations,
			Accesses:  stats.accesses,
		})
	}

	return trend, nil
}

// GetPasswordStrengths retrieves password strength analysis
func (r *EventSourcingRepository) GetPasswordStrengths(ctx context.Context, userID string) ([]models.PasswordStrength, error) {
	if userID == "" {
		return nil, ErrInvalidUserID
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	// Get all credential events for this user
	events, err := r.GetEventsByType(ctx, repositories.EventType(CredentialEvent))
	if err != nil {
		return nil, err
	}

	var strengths []models.PasswordStrength
	for _, event := range events {
		var data struct {
			UserID   string `json:"user_id"`
			Password string `json:"password"`
		}

		if err := json.Unmarshal(event.Data, &data); err != nil {
			continue
		}

		if data.UserID != userID {
			continue
		}

		strength := analyzePasswordStrength(data.Password)
		strengths = append(strengths, models.PasswordStrength{
			CredentialID: event.ID,
			Strength:     strength,
			Score:        calculatePasswordScore(data.Password),
		})
	}

	return strengths, nil
}

// GetReusedPasswords retrieves reused password analysis
func (r *EventSourcingRepository) GetReusedPasswords(ctx context.Context, userID string) ([]models.ReusedPassword, error) {
	if userID == "" {
		return nil, ErrInvalidUserID
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	// Get all credential events for this user
	events, err := r.GetEventsByType(ctx, repositories.EventType(CredentialEvent))
	if err != nil {
		return nil, err
	}

	passwordMap := make(map[string][]string)
	for _, event := range events {
		var data struct {
			UserID   string `json:"user_id"`
			Password string `json:"password"`
		}

		if err := json.Unmarshal(event.Data, &data); err != nil {
			continue
		}

		if data.UserID != userID {
			continue
		}

		passwordMap[data.Password] = append(passwordMap[data.Password], event.ID)
	}

	var reused []models.ReusedPassword
	for password, credentialIDs := range passwordMap {
		if len(credentialIDs) > 1 {
			reused = append(reused, models.ReusedPassword{
				Password:      password,
				CredentialIDs: credentialIDs,
			})
		}
	}

	return reused, nil
}

// GetBreachedCredentials retrieves breached credential analysis
func (r *EventSourcingRepository) GetBreachedCredentials(ctx context.Context, userID string) ([]models.BreachedCredential, error) {
	if userID == "" {
		return nil, ErrInvalidUserID
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	// Get all credential events for this user
	events, err := r.GetEventsByType(ctx, repositories.EventType(CredentialEvent))
	if err != nil {
		return nil, err
	}

	var breached []models.BreachedCredential
	for _, event := range events {
		var data struct {
			UserID   string `json:"user_id"`
			Password string `json:"password"`
		}

		if err := json.Unmarshal(event.Data, &data); err != nil {
			continue
		}

		if data.UserID != userID {
			continue
		}

		if isPasswordBreached(data.Password) {
			breached = append(breached, models.BreachedCredential{
				CredentialID: event.ID,
				Password:     data.Password,
			})
		}
	}

	return breached, nil
}

// GetOldPasswords retrieves old password analysis
func (r *EventSourcingRepository) GetOldPasswords(ctx context.Context, userID string, daysThreshold int) ([]models.OldPassword, error) {
	if userID == "" {
		return nil, ErrInvalidUserID
	}

	if daysThreshold <= 0 {
		return nil, ErrInvalidThreshold
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	// Get all credential events for this user
	events, err := r.GetEventsByType(ctx, repositories.EventType(CredentialEvent))
	if err != nil {
		return nil, err
	}

	var oldPasswords []models.OldPassword
	thresholdDate := time.Now().AddDate(0, 0, -daysThreshold)

	for _, event := range events {
		var data struct {
			UserID   string `json:"user_id"`
			Password string `json:"password"`
		}

		if err := json.Unmarshal(event.Data, &data); err != nil {
			continue
		}

		if data.UserID != userID {
			continue
		}

		if event.Timestamp.Before(thresholdDate) {
			oldPasswords = append(oldPasswords, models.OldPassword{
				CredentialID: event.ID,
				Password:     data.Password,
				Age:          int(time.Since(event.Timestamp).Hours() / 24),
			})
		}
	}

	return oldPasswords, nil
}

// GetCredentialAccesses returns all access events for a credential
func (r *EventSourcingRepository) GetCredentialAccesses(ctx context.Context, credentialID string, startDate, endDate time.Time) ([]models.CredentialAccess, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var accesses []models.CredentialAccess
	for _, event := range r.events {
		if event.Type == repositories.EventType(AccessEvent) {
			var access models.CredentialAccess
			if err := json.Unmarshal(event.Data, &access); err != nil {
				continue
			}
			if access.CredentialID == credentialID {
				accesses = append(accesses, access)
			}
		}
	}

	return accesses, nil
}

// Helper functions for password analysis
func analyzePasswordStrength(password string) string {
	// Implement password strength analysis logic
	return "strong" // Placeholder
}

func calculatePasswordScore(password string) int {
	// Implement password score calculation logic
	return 100 // Placeholder
}

func isPasswordBreached(password string) bool {
	// Implement password breach checking logic
	return false // Placeholder
}

func generateUUID() string {
	// Implement UUID generation logic
	return "uuid" // Placeholder
}

// Error definitions
var (
	ErrInvalidEventType    = errors.New("invalid event type")
	ErrInvalidEventData    = errors.New("invalid event data")
	ErrInvalidUserID       = errors.New("invalid user ID")
	ErrInvalidGroupID      = errors.New("invalid group ID")
	ErrInvalidCredentialID = errors.New("invalid credential ID")
	ErrInvalidTimeRange    = errors.New("invalid time range")
	ErrInvalidThreshold    = errors.New("invalid threshold")
)
