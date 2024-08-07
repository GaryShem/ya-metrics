package repository

import (
	"errors"

	"github.com/GaryShem/ya-metrics.git/internal/shared/storage/models"
)

var ErrMetricNotFound = errors.New("metric not found")

// Repository - generic metric storage interface.
type Repository interface {
	// UpdateGauge - updates gauge with specified key and value.
	UpdateGauge(metricName string, value float64) error
	// UpdateCounter - updates counter with specified key and value.
	UpdateCounter(metricName string, value int64) error
	// UpdateMetric - updates metric with specified incoming key and value.
	UpdateMetric(m *models.Metrics) error
	// UpdateMetricBatch - updates metric list with specified incoming key and value.
	UpdateMetricBatch(metrics []models.Metrics) ([]models.Metrics, error)
	// GetGauge - gets the gauge with specified key from storage.
	GetGauge(metricName string) (*models.Gauge, error)
	// GetCounter - gets the counter with specified key from storage.
	GetCounter(metricName string) (*models.Counter, error)
	// GetMetric - gets the metric with specified type and key from storage.
	GetMetric(m *models.Metrics) error
	// ListMetrics - returns a list with all metrics present in storage.
	ListMetrics() ([]models.Metrics, error)
	// GetGauges - returns a list with all gauges present in storage.
	GetGauges() (map[string]models.Gauge, error)
	// GetCounters - returns a list with all counters present in storage.
	GetCounters() (map[string]models.Counter, error)
	// Ping - heartbeat check for underlying storage.
	Ping() error
}
