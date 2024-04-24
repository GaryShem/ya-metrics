package memstorage

import (
	"fmt"

	"github.com/GaryShem/ya-metrics.git/internal/shared/storage/metrics"
)

func (ms *MemStorage) ResetCounter(metricName string) error {
	if counter, ok := ms.CounterMetrics[metricName]; ok {
		counter.Value = 0
	} else {
		return fmt.Errorf("%w: %v", ErrMetricNotFound, metricName)
	}
	return nil
}

func (ms *MemStorage) GetCounters() map[string]*metrics.Counter {
	result := make(map[string]*metrics.Counter)
	for k, v := range ms.CounterMetrics {
		result[k] = metrics.CopyCounter(*v)
	}
	return result
}

func (ms *MemStorage) GetCounter(metricName string) (*metrics.Counter, error) {
	value, ok := ms.CounterMetrics[metricName]
	if !ok {
		return nil, fmt.Errorf("%w: %v", ErrMetricNotFound, metricName)
	}
	return metrics.CopyCounter(*value), nil
}

func (ms *MemStorage) UpdateCounter(metricName string, delta int64) {
	currentValue, ok := ms.CounterMetrics[metricName]
	if !ok {
		ms.CounterMetrics[metricName] = metrics.NewCounter(metricName, delta)
	} else {
		currentValue.Update(delta)
	}
}
