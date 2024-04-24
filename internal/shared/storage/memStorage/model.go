package memstorage

import (
	"errors"

	"github.com/GaryShem/ya-metrics.git/internal/shared/storage"
	"github.com/GaryShem/ya-metrics.git/internal/shared/storage/metrics"
)

var ErrMetricNotFound = errors.New("metric not found")

type MemStorage struct {
	GaugeMetrics   map[string]*metrics.Gauge   `json:"gaugeMetrics"`
	CounterMetrics map[string]*metrics.Counter `json:"counterMetrics"`
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		GaugeMetrics:   make(map[string]*metrics.Gauge),
		CounterMetrics: make(map[string]*metrics.Counter),
	}
}

var _ storage.Repository = &MemStorage{}
