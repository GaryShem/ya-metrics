package storage

import (
	"fmt"

	"github.com/GaryShem/ya-metrics.git/internal/shared/storage/metrics"
)

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

func (ms *MemStorage) ResetCounter(metricName string) error {
	if counter, ok := ms.CounterMetrics[metricName]; ok {
		counter.Value = 0
	} else {
		return fmt.Errorf("metric counter %s not present in storage", metricName)
	}
	return nil
}

func (ms *MemStorage) GetGauges() map[string]*metrics.Gauge {
	result := make(map[string]*metrics.Gauge)
	for k, v := range ms.GaugeMetrics {
		result[k] = metrics.CopyGauge(*v)
	}
	return result
}

func (ms *MemStorage) GetCounters() map[string]*metrics.Counter {
	result := make(map[string]*metrics.Counter)
	for k, v := range ms.CounterMetrics {
		result[k] = metrics.CopyCounter(*v)
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

func (ms *MemStorage) UpdateCounter(metricName string, value int64) {
	currentValue, ok := ms.CounterMetrics[metricName]
	if !ok {
		ms.CounterMetrics[metricName] = metrics.NewCounter(metricName, value)
	} else {
		currentValue.Update(value)
	}
}

func (ms *MemStorage) GetGauge(metricName string) (*metrics.Gauge, error) {
	value, ok := ms.GaugeMetrics[metricName]
	if !ok {
		return nil, fmt.Errorf("no value for metric %s", metricName)
	}
	return metrics.CopyGauge(*value), nil
}

func (ms *MemStorage) GetCounter(metricName string) (*metrics.Counter, error) {
	value, ok := ms.CounterMetrics[metricName]
	if !ok {
		return nil, fmt.Errorf("no value for metric %s", metricName)
	}
	return metrics.CopyCounter(*value), nil
}

var _ Repository = &MemStorage{}
