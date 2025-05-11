package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	esdb "github.com/kurrent-io/KurrentDB-Client-Go/kurrentdb"
	"github.com/polypass/polypass-microservices/statistics/domain/models"
	"github.com/polypass/polypass-microservices/statistics/domain/repositories"
)

// EventStoreDBEventRepository implements the EventRepository interface using EventStoreDB
type EventStoreDBEventRepository struct {
	client *esdb.Client
}

// NewEventStoreDBEventRepository creates a new instance of EventStoreDBEventRepository
func NewEventStoreDBEventRepository(client *esdb.Client) repositories.EventRepository {
	return &EventStoreDBEventRepository{
		client: client,
	}
}

// StoreEvent persists an event to EventStoreDB
func (r *EventStoreDBEventRepository) StoreEvent(ctx context.Context, event models.Event) error {
	// Convert the event to JSON
	eventData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Create EventStoreDB event data
	eventDataBytes := esdb.EventData{
		ContentType: esdb.ContentTypeJson,
		EventType:   event.Type,
		Data:        eventData,
	}

	// Create options for appending the event
	opts := esdb.AppendToStreamOptions{}

	// Append the event to the stream
	streamName := fmt.Sprintf("events-%s", event.Type)
	_, err = r.client.AppendToStream(ctx, streamName, opts, eventDataBytes)
	if err != nil {
		return fmt.Errorf("failed to append event to stream: %w", err)
	}

	return nil
}

// GetEventByID retrieves an event by its ID from EventStoreDB
func (r *EventStoreDBEventRepository) GetEventByID(ctx context.Context, id string) (models.Event, error) {
	var event models.Event

	// Read all events from all streams
	opts := esdb.ReadAllOptions{
		Direction:      esdb.Forwards,
		ResolveLinkTos: true,
	}

	stream, err := r.client.ReadAll(ctx, opts, uint64(18446744073709551615))
	if err != nil {
		return event, fmt.Errorf("failed to read all events: %w", err)
	}
	defer stream.Close()

	// Iterate through all events to find the one with the matching ID
	for {
		resolvedEvent, err := stream.Recv()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return event, fmt.Errorf("failed to receive event: %w", err)
		}

		// Skip system events
		if resolvedEvent.Event == nil || resolvedEvent.Event.EventType == "$" {
			continue
		}

		// Unmarshal the event data
		var e models.Event
		if err := json.Unmarshal(resolvedEvent.Event.Data, &e); err != nil {
			continue
		}

		// Check if this is the event we're looking for
		if e.ID == id {
			return e, nil
		}
	}

	return event, fmt.Errorf("event with ID %s not found", id)
}

// GetEventsByType retrieves events of a specific type from EventStoreDB
func (r *EventStoreDBEventRepository) GetEventsByType(ctx context.Context, eventType models.EventType) ([]models.Event, error) {
	streamName := fmt.Sprintf("events-%s", string(eventType))
	return r.readEventsFromStream(ctx, streamName)
}

// GetEventsByTimeRange retrieves events within a specific time range from EventStoreDB
func (r *EventStoreDBEventRepository) GetEventsByTimeRange(ctx context.Context, start, end time.Time) ([]models.Event, error) {
	// Read all events from all streams
	opts := esdb.ReadAllOptions{
		Direction:      esdb.Forwards,
		ResolveLinkTos: true,
	}

	stream, err := r.client.ReadAll(ctx, opts, uint64(18446744073709551615))
	if err != nil {
		return nil, fmt.Errorf("failed to read all events: %w", err)
	}
	defer stream.Close()

	var events []models.Event

	// Iterate through all events to find those within the time range
	for {
		resolvedEvent, err := stream.Recv()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, fmt.Errorf("failed to receive event: %w", err)
		}

		// Skip system events
		if resolvedEvent.Event == nil || resolvedEvent.Event.EventType == "$" {
			continue
		}

		// Unmarshal the event data
		var event models.Event
		if err := json.Unmarshal(resolvedEvent.Event.Data, &event); err != nil {
			continue
		}

		// Check if the event is within the time range
		if (event.Timestamp.Equal(start) || event.Timestamp.After(start)) &&
			(event.Timestamp.Equal(end) || event.Timestamp.Before(end)) {
			events = append(events, event)
		}
	}

	return events, nil
}

// GetEventsByTypeAndTimeRange retrieves events of a specific type within a time range from EventStoreDB
func (r *EventStoreDBEventRepository) GetEventsByTypeAndTimeRange(ctx context.Context, eventType models.EventType, start, end time.Time) ([]models.Event, error) {
	// Get all events of the specified type
	eventsOfType, err := r.GetEventsByType(ctx, eventType)
	if err != nil {
		return nil, err
	}

	// Filter events by time range
	var events []models.Event
	for _, event := range eventsOfType {
		if (event.Timestamp.Equal(start) || event.Timestamp.After(start)) &&
			(event.Timestamp.Equal(end) || event.Timestamp.Before(end)) {
			events = append(events, event)
		}
	}

	return events, nil
}

// GetLatestEvents retrieves the most recent events from EventStoreDB, limited by count
func (r *EventStoreDBEventRepository) GetLatestEvents(ctx context.Context, limit int) ([]models.Event, error) {
	// Read all events from all streams in reverse order
	opts := esdb.ReadAllOptions{
		Direction:      esdb.Backwards,
		ResolveLinkTos: true,
	}

	stream, err := r.client.ReadAll(ctx, opts, uint64(limit))
	if err != nil {
		return nil, fmt.Errorf("failed to read all events: %w", err)
	}
	defer stream.Close()

	var events []models.Event
	count := 0

	// Iterate through events to get the latest ones
	for {
		if count >= limit {
			break
		}

		resolvedEvent, err := stream.Recv()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, fmt.Errorf("failed to receive event: %w", err)
		}

		// Skip system events
		if resolvedEvent.Event == nil || resolvedEvent.Event.EventType == "$" {
			continue
		}

		// Unmarshal the event data
		var event models.Event
		if err := json.Unmarshal(resolvedEvent.Event.Data, &event); err != nil {
			continue
		}

		events = append(events, event)
		count++
	}

	return events, nil
}

// readEventsFromStream is a helper method to read events from a specific stream
func (r *EventStoreDBEventRepository) readEventsFromStream(ctx context.Context, streamName string) ([]models.Event, error) {
	// Read all events from the stream
	opts := esdb.ReadStreamOptions{
		Direction: esdb.Forwards,
		From:      esdb.Start{},
	}

	stream, err := r.client.ReadStream(ctx, streamName, opts, 1000)
	if err != nil {
		return nil, fmt.Errorf("failed to read stream: %w", err)
	}
	defer stream.Close()

	var events []models.Event

	// Iterate through all events in the stream
	for {
		resolvedEvent, err := stream.Recv()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, fmt.Errorf("failed to receive event: %w", err)
		}

		// Unmarshal the event data
		var event models.Event
		if err := json.Unmarshal(resolvedEvent.Event.Data, &event); err != nil {
			continue
		}

		events = append(events, event)
	}

	return events, nil
}
