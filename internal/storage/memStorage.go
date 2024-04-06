package storage

import "fmt"

type MemStorage struct {
	GaugeMetrics   map[string]float64 `json:"gaugeMetrics"`
	CounterMetrics map[string]int64   `json:"counterMetrics"`
}

type repository interface {
	UpdateGauge(metricName string, value float64)
	UpdateCounter(metricName string, value int64)
	GetGauge(metricName string) (float64, error)
	GetCounter(metricName string) (int64, error)
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		GaugeMetrics:   make(map[string]float64),
		CounterMetrics: make(map[string]int64),
	}
}

func (ms *MemStorage) UpdateGauge(metricName string, value float64) {
	ms.GaugeMetrics[metricName] = value
}

func (ms *MemStorage) UpdateCounter(metricName string, value int64) {
	currentValue := ms.CounterMetrics[metricName]
	ms.CounterMetrics[metricName] = currentValue + value
}

func (ms *MemStorage) GetGauge(metricName string) (float64, error) {
	value, ok := ms.GaugeMetrics[metricName]
	if !ok {
		return 0, fmt.Errorf("no value for metric %s", metricName)
	}
	return value, nil
}

func (ms *MemStorage) GetCounter(metricName string) (int64, error) {
	value, ok := ms.CounterMetrics[metricName]
	if !ok {
		return 0, fmt.Errorf("no value for metric %s", metricName)
	}
	return value, nil
}

var _ repository = &MemStorage{}
