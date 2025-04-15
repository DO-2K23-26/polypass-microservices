package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/DO-2K23-26/polypass-microservices/statistics/domain/models"
)

// EventType represents the type of an event
type EventType string

const (
	CredentialEvent EventType = "credential"
	AccessEvent     EventType = "access"
)

// Event represents a domain event
type Event struct {
	ID        string
	Type      EventType
	Data      []byte
	Timestamp time.Time
}

// EventSourcingRepository defines the interface for event sourcing operations
type EventSourcingRepository interface {
	// Setup initializes the repository
	Setup() error

	// Shutdown cleans up resources
	Shutdown() error

	// StoreEvent stores an event
	StoreEvent(ctx context.Context, event Event) error

	// GetEventsByType retrieves events of a specific type
	GetEventsByType(ctx context.Context, eventType EventType) ([]Event, error)

	// GetEventsByTimeRange retrieves events within a time range
	GetEventsByTimeRange(ctx context.Context, start, end time.Time) ([]Event, error)

	// GetUserMetrics retrieves user metrics by replaying events
	GetUserMetrics(ctx context.Context, userID string) (*models.UserMetrics, error)

	// GetGroupMetrics retrieves group metrics by replaying events
	GetGroupMetrics(ctx context.Context, groupID string) (*models.GroupMetrics, error)

	// GetCredentialAccesses retrieves credential access records
	GetCredentialAccesses(ctx context.Context, credentialID string, startDate, endDate time.Time) ([]models.CredentialAccess, error)

	// GetSharedCredentialStats retrieves shared credential statistics
	GetSharedCredentialStats(ctx context.Context, credentialID string) (*models.SharedCredentialStats, error)

	// GetCredentialTrends retrieves credential creation and usage trends
	GetCredentialTrends(ctx context.Context, startDate, endDate time.Time) (*models.CredentialTrend, error)

	// GetPasswordStrengths retrieves password strength analysis
	GetPasswordStrengths(ctx context.Context, userID string) ([]models.PasswordStrength, error)

	// GetReusedPasswords retrieves reused password analysis
	GetReusedPasswords(ctx context.Context, userID string) ([]models.ReusedPassword, error)

	// GetBreachedCredentials retrieves breached credential analysis
	GetBreachedCredentials(ctx context.Context, userID string) ([]models.BreachedCredential, error)

	// GetOldPasswords retrieves old password analysis
	GetOldPasswords(ctx context.Context, userID string, daysThreshold int) ([]models.OldPassword, error)
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
