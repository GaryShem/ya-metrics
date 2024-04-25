package memstorage

import (
	"fmt"

	"github.com/GaryShem/ya-metrics.git/internal/shared/storage/models"
)

func (ms *MemStorage) ResetCounter(metricName string) error {
	if counter, ok := ms.CounterMetrics[metricName]; ok {
		counter.Value = 0
	} else {
		return fmt.Errorf("%w: %v", ErrMetricNotFound, metricName)
	}
	return nil
}

func (ms *MemStorage) GetCounters() map[string]*models.Counter {
	result := make(map[string]*models.Counter)
	for k, v := range ms.CounterMetrics {
		result[k] = models.CopyCounter(*v)
	}
	return result
}

func (ms *MemStorage) GetCounter(metricName string) (*models.Counter, error) {
	value, ok := ms.CounterMetrics[metricName]
	if !ok {
		return nil, fmt.Errorf("%w: %v", ErrMetricNotFound, metricName)
	}
	return models.CopyCounter(*value), nil
}

func (ms *MemStorage) UpdateCounter(metricName string, delta int64) {
	currentValue, ok := ms.CounterMetrics[metricName]
	if !ok {
		ms.CounterMetrics[metricName] = models.NewCounter(metricName, delta)
	} else {
		currentValue.Update(delta)
	}
}
