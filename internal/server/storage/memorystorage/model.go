package memorystorage

import (
	"sync"
	"time"

	"github.com/GaryShem/ya-metrics.git/internal/server/storage/repository"
	"github.com/GaryShem/ya-metrics.git/internal/shared/storage/models"
)

type MemStorage struct {
	GaugeMetrics   map[string]*models.Gauge   `json:"gaugeMetrics"`
	CounterMetrics map[string]*models.Counter `json:"counterMetrics"`
	mu             sync.RWMutex
	LastChangeTime time.Time `json:"lastChangeTime"`
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		GaugeMetrics:   make(map[string]*models.Gauge),
		CounterMetrics: make(map[string]*models.Counter),
	}
}

var _ repository.Repository = &MemStorage{}
