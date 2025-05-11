package repositories

import (
	"context"
	"time"

	"github.com/polypass/polypass-microservices/statistics/domain/models"
)

// EventRepository defines the interface for event storage operations
type EventRepository interface {
	// StoreEvent persists an event to the storage
	StoreEvent(ctx context.Context, event models.Event) error

	// GetEventByID retrieves an event by its ID
	GetEventByID(ctx context.Context, id string) (models.Event, error)

	// GetEventsByType retrieves events of a specific type
	GetEventsByType(ctx context.Context, eventType models.EventType) ([]models.Event, error)

	// GetEventsByTimeRange retrieves events within a specific time range
	GetEventsByTimeRange(ctx context.Context, start, end time.Time) ([]models.Event, error)

	// GetEventsByTypeAndTimeRange retrieves events of a specific type within a time range
	GetEventsByTypeAndTimeRange(ctx context.Context, eventType models.EventType, start, end time.Time) ([]models.Event, error)

	// GetLatestEvents retrieves the most recent events, limited by count
	GetLatestEvents(ctx context.Context, limit int) ([]models.Event, error)
}
