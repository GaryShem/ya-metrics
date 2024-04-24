package memstorage

import (
	"fmt"

	"github.com/GaryShem/ya-metrics.git/internal/shared/storage/metrics"
)

func (ms *MemStorage) GetGauges() map[string]*metrics.Gauge {
	result := make(map[string]*metrics.Gauge)
	for k, v := range ms.GaugeMetrics {
		result[k] = metrics.CopyGauge(*v)
	}
	return result
}

func (ms *MemStorage) UpdateGauge(metricName string, value float64) {
	currentValue, ok := ms.GaugeMetrics[metricName]
	if !ok {
		ms.GaugeMetrics[metricName] = metrics.NewGauge(metricName, value)
	} else {
		currentValue.Update(value)
	}
}

func (ms *MemStorage) GetGauge(metricName string) (*metrics.Gauge, error) {
	value, ok := ms.GaugeMetrics[metricName]
	if !ok {
		return nil, fmt.Errorf("%w: %v", ErrMetricNotFound, metricName)
	}
	return metrics.CopyGauge(*value), nil
}
