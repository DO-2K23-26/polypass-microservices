package repositories

import (
	"context"
	"time"

	"github.com/polypass/polypass-microservices/statistics/domain/models"
)

// MetricRepository defines the interface for metric storage operations
type MetricRepository interface {
	// StoreMetric persists a metric to the storage
	StoreMetric(ctx context.Context, metric models.Metric) error

	// GetMetricByID retrieves a metric by its ID
	GetMetricByID(ctx context.Context, id string) (models.Metric, error)

	// GetMetricsByName retrieves metrics with a specific name
	GetMetricsByName(ctx context.Context, name string) ([]models.Metric, error)

	// GetMetricsByTimeRange retrieves metrics within a specific time range
	GetMetricsByTimeRange(ctx context.Context, start, end time.Time) ([]models.Metric, error)

	// GetMetricsByCategory retrieves metrics of a specific category
	GetMetricsByCategory(ctx context.Context, category models.MetricCategory) ([]models.Metric, error)

	// GetMetricsByCategoryAndTimeRange retrieves metrics of a specific category within a time range
	GetMetricsByCategoryAndTimeRange(ctx context.Context, category models.MetricCategory, start, end time.Time) ([]models.Metric, error)

	// GetLatestMetrics retrieves the most recent metrics, limited by count
	GetLatestMetrics(ctx context.Context, limit int) ([]models.Metric, error)
}
