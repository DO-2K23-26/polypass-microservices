package services

import (
	"context"
	"time"

	"github.com/polypass/polypass-microservices/statistics/domain/models"
	"github.com/polypass/polypass-microservices/statistics/domain/repositories"
)

// MetricService handles operations related to metrics
type MetricService struct {
	metricRepository repositories.MetricRepository
	eventService     *EventService
	calculators      []models.MetricCalculator
}

// NewMetricService creates a new instance of MetricService
func NewMetricService(
	metricRepo repositories.MetricRepository,
	eventService *EventService,
	calculators []models.MetricCalculator,
) *MetricService {
	return &MetricService{
		metricRepository: metricRepo,
		eventService:     eventService,
		calculators:      calculators,
	}
}

// RegisterCalculator adds a new metric calculator to the service
func (s *MetricService) RegisterCalculator(calculator models.MetricCalculator) {
	s.calculators = append(s.calculators, calculator)
}

// CalculateMetrics calculates all metrics using registered calculators
func (s *MetricService) CalculateMetrics(ctx context.Context, timeRange time.Duration) ([]models.Metric, error) {
	// Get events for the specified time range
	end := time.Now()
	start := end.Add(-timeRange)

	events, err := s.eventService.GetEventsByTimeRange(ctx, start, end)
	if err != nil {
		return nil, err
	}

	var metrics []models.Metric

	// Calculate metrics using each registered calculator
	for _, calculator := range s.calculators {
		metric, err := calculator.Calculate(events)
		if err != nil {
			// Log error but continue with other calculators
			continue
		}

		// Store the calculated metric
		err = s.metricRepository.StoreMetric(ctx, metric)
		if err != nil {
			// Log error but continue with other metrics
			continue
		}

		metrics = append(metrics, metric)
	}

	return metrics, nil
}

// GetMetricByID retrieves a metric by its ID
func (s *MetricService) GetMetricByID(ctx context.Context, id string) (models.Metric, error) {
	return s.metricRepository.GetMetricByID(ctx, id)
}

// GetMetricsByName retrieves metrics with a specific name
func (s *MetricService) GetMetricsByName(ctx context.Context, name string) ([]models.Metric, error) {
	return s.metricRepository.GetMetricsByName(ctx, name)
}

// GetMetricsByCategory retrieves metrics of a specific category
func (s *MetricService) GetMetricsByCategory(ctx context.Context, category models.MetricCategory) ([]models.Metric, error) {
	return s.metricRepository.GetMetricsByCategory(ctx, category)
}

// GetLatestMetrics retrieves the most recent metrics, limited by count
func (s *MetricService) GetLatestMetrics(ctx context.Context, limit int) ([]models.Metric, error) {
	return s.metricRepository.GetLatestMetrics(ctx, limit)
}
