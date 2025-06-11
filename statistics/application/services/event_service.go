package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/polypass/polypass-microservices/statistics/domain/models"
	"github.com/polypass/polypass-microservices/statistics/domain/repositories"
)

// EventService handles operations related to events
type EventService struct {
	eventRepository repositories.EventRepository
}

// NewEventService creates a new instance of EventService
func NewEventService(eventRepo repositories.EventRepository) *EventService {
	return &EventService{
		eventRepository: eventRepo,
	}
}

// StoreEvent stores an event in the repository
func (s *EventService) StoreEvent(ctx context.Context, event models.Event) error {
	return s.eventRepository.StoreEvent(ctx, event)
}

// GetEventByID retrieves an event by its ID
func (s *EventService) GetEventByID(ctx context.Context, id string) (models.Event, error) {
	return s.eventRepository.GetEventByID(ctx, id)
}

// GetEventsByType retrieves events of a specific type
func (s *EventService) GetEventsByType(ctx context.Context, eventType models.EventType) ([]models.Event, error) {
	return s.eventRepository.GetEventsByType(ctx, eventType)
}

// GetEventsByTimeRange retrieves events within a specific time range
func (s *EventService) GetEventsByTimeRange(ctx context.Context, start, end time.Time) ([]models.Event, error) {
	return s.eventRepository.GetEventsByTimeRange(ctx, start, end)
}

// GetEventsByTypeAndTimeRange retrieves events of a specific type within a time range
func (s *EventService) GetEventsByTypeAndTimeRange(ctx context.Context, eventType models.EventType, start, end time.Time) ([]models.Event, error) {
	return s.eventRepository.GetEventsByTypeAndTimeRange(ctx, eventType, start, end)
}

// GetLatestEvents retrieves the most recent events, limited by count
func (s *EventService) GetLatestEvents(ctx context.Context, limit int) ([]models.Event, error) {
	return s.eventRepository.GetLatestEvents(ctx, limit)
}

func (s *EventService) UsedPasswordCount(ctx context.Context, req struct{UserID string `json:"userId"`; PasswordID string `json:"passwordId"`}) (error) {
	events, err := s.eventRepository.GetEventsByTypeAndTimeRange(ctx, models.PasswordUsedEvent, time.Now().Add(-time.Hour), time.Now())
	if err != nil {
		return err
	}
	for _, event := range events {
		dataMap, ok := event.Data.(map[string]interface{})
		if !ok {
			continue
		} 
		if dataMap["passwordId"] == req.PasswordID {
			return nil
		}
	}
	newEvent := models.Event{
		ID:        uuid.New().String(),
		Type:      string(models.PasswordUsedEvent),
		Data:      map[string]interface{}{"passwordId": req.PasswordID},
		Timestamp: time.Now(),
	}
	if err := s.eventRepository.StoreEvent(ctx, newEvent); err != nil {
		return err
	}
	return nil
}
