package repository

import (
	"errors"

	"github.com/GaryShem/ya-metrics.git/internal/shared/storage/models"
)

var ErrMetricNotFound = errors.New("metric not found")

type Repository interface {
	UpdateGauge(metricName string, value float64) error
	UpdateCounter(metricName string, value int64) error
	UpdateMetric(m *models.Metrics) error
	UpdateMetricBatch(metrics []models.Metrics) ([]models.Metrics, error)
	GetGauge(metricName string) (*models.Gauge, error)
	GetCounter(metricName string) (*models.Counter, error)
	GetMetric(m *models.Metrics) error
	ListMetrics() ([]models.Metrics, error)
	GetGauges() (map[string]models.Gauge, error)
	GetCounters() (map[string]models.Counter, error)
	Ping() error
}
