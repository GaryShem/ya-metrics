package mem_storage

import (
	"errors"

	"github.com/GaryShem/ya-metrics.git/internal/shared/storage/models"
)

var ErrMetricNotFound = errors.New("metric not found")

type MemStorage struct {
	GaugeMetrics   map[string]*models.Gauge   `json:"gaugeMetrics"`
	CounterMetrics map[string]*models.Counter `json:"counterMetrics"`
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		GaugeMetrics:   make(map[string]*models.Gauge),
		CounterMetrics: make(map[string]*models.Counter),
	}
}

var _ models.Repository = &MemStorage{}
