package storage

import "github.com/GaryShem/ya-metrics.git/internal/shared/storage/metrics"

type Repository interface {
	UpdateGauge(metricName string, value float64)
	UpdateCounter(metricName string, value int64)
	GetGauge(metricName string) (*metrics.Gauge, error)
	GetCounter(metricName string) (*metrics.Counter, error)
	GetGauges() map[string]*metrics.Gauge
	GetCounters() map[string]*metrics.Counter

	ResetCounter(metricName string) error
}