package memorystorage

import (
	"fmt"
	"time"

	"github.com/GaryShem/ya-metrics.git/internal/shared/storage/models"
)

func (ms *MemStorage) GetCounters() map[string]*models.Counter {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	result := make(map[string]*models.Counter)
	for k, v := range ms.CounterMetrics {
		result[k] = models.CopyCounter(*v)
	}
	return result
}

func (ms *MemStorage) GetCounter(metricName string) (*models.Counter, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	value, ok := ms.CounterMetrics[metricName]
	if !ok {
		return nil, fmt.Errorf("%w: %v", ErrMetricNotFound, metricName)
	}
	return models.CopyCounter(*value), nil
}

func (ms *MemStorage) UpdateCounter(metricName string, delta int64) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	currentValue, ok := ms.CounterMetrics[metricName]
	if !ok {
		ms.CounterMetrics[metricName] = models.NewCounter(metricName, delta)
	} else {
		currentValue.Update(delta)
	}
	ms.LastChangeTime = time.Now()
}
