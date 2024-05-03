package memorystorage

import (
	"fmt"
	"time"

	"github.com/GaryShem/ya-metrics.git/internal/shared/storage/models"
)

func (ms *MemStorage) GetGauges() map[string]*models.Gauge {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	result := make(map[string]*models.Gauge)
	for k, v := range ms.GaugeMetrics {
		result[k] = models.CopyGauge(*v)
	}
	return result
}

func (ms *MemStorage) GetGauge(metricName string) (*models.Gauge, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	value, ok := ms.GaugeMetrics[metricName]
	if !ok {
		return nil, fmt.Errorf("%w: %v", ErrMetricNotFound, metricName)
	}
	return models.CopyGauge(*value), nil
}

func (ms *MemStorage) UpdateGauge(metricName string, value float64) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	currentValue, ok := ms.GaugeMetrics[metricName]
	if !ok {
		ms.GaugeMetrics[metricName] = models.NewGauge(metricName, value)
	} else {
		currentValue.Update(value)
	}
	ms.LastChangeTime = time.Now()
}
