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

// EventStoreDBMetricRepository implements the MetricRepository interface using EventStoreDB
type EventStoreDBMetricRepository struct {
	client *esdb.Client
}

// NewEventStoreDBMetricRepository creates a new instance of EventStoreDBMetricRepository
func NewEventStoreDBMetricRepository(client *esdb.Client) repositories.MetricRepository {
	return &EventStoreDBMetricRepository{
		client: client,
	}
}

// StoreMetric persists a metric to EventStoreDB
func (r *EventStoreDBMetricRepository) StoreMetric(ctx context.Context, metric models.Metric) error {
	// Convert the metric to JSON
	metricData, err := json.Marshal(metric)
	if err != nil {
		return fmt.Errorf("failed to marshal metric: %w", err)
	}

	// Create EventStoreDB event data
	eventDataBytes := esdb.EventData{
		ContentType: esdb.ContentTypeJson,
		EventType:   "METRIC",
		Data:        metricData,
	}

	// Create options for appending the event
	opts := esdb.AppendToStreamOptions{}

	// Append the metric to the stream
	streamName := fmt.Sprintf("metrics-%s", metric.Name)
	_, err = r.client.AppendToStream(ctx, streamName, opts, eventDataBytes)
	if err != nil {
		return fmt.Errorf("failed to append metric to stream: %w", err)
	}

	return nil
}

// GetMetricByID retrieves a metric by its ID from EventStoreDB
func (r *EventStoreDBMetricRepository) GetMetricByID(ctx context.Context, id string) (models.Metric, error) {
	var metric models.Metric

	// Read all events from all streams
	opts := esdb.ReadAllOptions{
		Direction:      esdb.Forwards,
		ResolveLinkTos: true,
	}

	stream, err := r.client.ReadAll(ctx, opts, uint64(18446744073709551615))
	if err != nil {
		return metric, fmt.Errorf("failed to read all events: %w", err)
	}
	defer stream.Close()

	// Iterate through all events to find the one with the matching ID
	for {
		resolvedEvent, err := stream.Recv()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return metric, fmt.Errorf("failed to receive event: %w", err)
		}

		// Skip system events and non-metric events
		if resolvedEvent.Event == nil || resolvedEvent.Event.EventType != "METRIC" {
			continue
		}

		// Unmarshal the metric data
		var m models.Metric
		if err := json.Unmarshal(resolvedEvent.Event.Data, &m); err != nil {
			continue
		}

		// Check if this is the metric we're looking for
		if m.ID == id {
			return m, nil
		}
	}

	return metric, fmt.Errorf("metric with ID %s not found", id)
}

// GetMetricsByName retrieves metrics with a specific name from EventStoreDB
func (r *EventStoreDBMetricRepository) GetMetricsByName(ctx context.Context, name string) ([]models.Metric, error) {
	streamName := fmt.Sprintf("metrics-%s", name)
	return r.readMetricsFromStream(ctx, streamName)
}

// GetMetricsByTimeRange retrieves metrics within a specific time range from EventStoreDB
func (r *EventStoreDBMetricRepository) GetMetricsByTimeRange(ctx context.Context, start, end time.Time) ([]models.Metric, error) {
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

	var metrics []models.Metric

	// Iterate through all events to find those within the time range
	for {
		resolvedEvent, err := stream.Recv()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, fmt.Errorf("failed to receive event: %w", err)
		}

		// Skip system events and non-metric events
		if resolvedEvent.Event == nil || resolvedEvent.Event.EventType != "METRIC" {
			continue
		}

		// Unmarshal the metric data
		var metric models.Metric
		if err := json.Unmarshal(resolvedEvent.Event.Data, &metric); err != nil {
			continue
		}

		// Check if the metric is within the time range
		if (metric.Timestamp.Equal(start) || metric.Timestamp.After(start)) &&
			(metric.Timestamp.Equal(end) || metric.Timestamp.Before(end)) {
			metrics = append(metrics, metric)
		}
	}

	return metrics, nil
}

// GetMetricsByCategory retrieves metrics of a specific category from EventStoreDB
func (r *EventStoreDBMetricRepository) GetMetricsByCategory(ctx context.Context, category models.MetricCategory) ([]models.Metric, error) {
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

	var metrics []models.Metric

	// Iterate through all events to find those with the specified category
	for {
		resolvedEvent, err := stream.Recv()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, fmt.Errorf("failed to receive event: %w", err)
		}

		// Skip system events and non-metric events
		if resolvedEvent.Event == nil || resolvedEvent.Event.EventType != "METRIC" {
			continue
		}

		// Unmarshal the metric data
		var metric models.Metric
		if err := json.Unmarshal(resolvedEvent.Event.Data, &metric); err != nil {
			continue
		}

		// Check if the metric has the specified category
		for _, tag := range metric.Tags {
			if tag == string(category) {
				metrics = append(metrics, metric)
				break
			}
		}
	}

	return metrics, nil
}

// GetMetricsByCategoryAndTimeRange retrieves metrics of a specific category within a time range from EventStoreDB
func (r *EventStoreDBMetricRepository) GetMetricsByCategoryAndTimeRange(ctx context.Context, category models.MetricCategory, start, end time.Time) ([]models.Metric, error) {
	// Get all metrics of the specified category
	metricsOfCategory, err := r.GetMetricsByCategory(ctx, category)
	if err != nil {
		return nil, err
	}

	// Filter metrics by time range
	var metrics []models.Metric
	for _, metric := range metricsOfCategory {
		if (metric.Timestamp.Equal(start) || metric.Timestamp.After(start)) &&
			(metric.Timestamp.Equal(end) || metric.Timestamp.Before(end)) {
			metrics = append(metrics, metric)
		}
	}

	return metrics, nil
}

// GetLatestMetrics retrieves the most recent metrics from EventStoreDB, limited by count
func (r *EventStoreDBMetricRepository) GetLatestMetrics(ctx context.Context, limit int) ([]models.Metric, error) {
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

	var metrics []models.Metric
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

		// Skip system events and non-metric events
		if resolvedEvent.Event == nil || resolvedEvent.Event.EventType != "METRIC" {
			continue
		}

		// Unmarshal the metric data
		var metric models.Metric
		if err := json.Unmarshal(resolvedEvent.Event.Data, &metric); err != nil {
			continue
		}

		metrics = append(metrics, metric)
		count++
	}

	return metrics, nil
}

// readMetricsFromStream is a helper method to read metrics from a specific stream
func (r *EventStoreDBMetricRepository) readMetricsFromStream(ctx context.Context, streamName string) ([]models.Metric, error) {
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

	var metrics []models.Metric

	// Iterate through all events in the stream
	for {
		resolvedEvent, err := stream.Recv()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, fmt.Errorf("failed to receive event: %w", err)
		}

		// Unmarshal the metric data
		var metric models.Metric
		if err := json.Unmarshal(resolvedEvent.Event.Data, &metric); err != nil {
			continue
		}

		metrics = append(metrics, metric)
	}

	return metrics, nil
}
