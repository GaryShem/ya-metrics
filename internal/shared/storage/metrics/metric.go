package metrics

import (
	"errors"
	"fmt"
	"slices"
)

var ErrInvalidMetricID = errors.New("metric ID is empty")
var ErrInvalidMetricType = errors.New("invalid metric type")
var ErrInvalidMetricValue = errors.New("empty metric value")

type MetricType string

const (
	TypeGauge   MetricType = "gauge"
	TypeCounter MetricType = "counter"
)

func GetSupportedMetricTypes() []MetricType {
	return []MetricType{TypeGauge, TypeCounter}
}

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (m *Metrics) ValidateUpdate() error {
	if m.ID == "" {
		return ErrInvalidMetricID
	}
	switch MetricType(m.MType) {
	case TypeCounter:
		if m.Delta == nil {
			return fmt.Errorf("%w: %s", ErrInvalidMetricValue, m.MType)
		}
	case TypeGauge:
		if m.Value == nil {
			return fmt.Errorf("%w: %s", ErrInvalidMetricValue, m.MType)
		}
	default:
		return fmt.Errorf("%w: %s", ErrInvalidMetricType, m.MType)
	}
	return nil
}

func (m *Metrics) ValidateGet() error {
	if m.ID == "" {
		return ErrInvalidMetricID
	}
	if !slices.Contains(GetSupportedMetricTypes(), MetricType(m.MType)) {
		return fmt.Errorf("%w: %s", ErrInvalidMetricType, m.MType)
	}
	return nil
}
