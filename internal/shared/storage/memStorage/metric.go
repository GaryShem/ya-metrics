package memstorage

import (
	"github.com/GaryShem/ya-metrics.git/internal/shared/storage/models"
)

func (ms *MemStorage) UpdateMetric(m *models.Metrics) error {
	if err := m.ValidateUpdate(); err != nil {
		return err
	}
	switch models.MetricType(m.MType) {
	case models.TypeGauge:
		ms.UpdateGauge(m.ID, *m.Value)
		v, err := ms.GetGauge(m.ID)
		if err != nil {
			return err
		}
		m.Value = &v.Value
		return nil
	case models.TypeCounter:
		ms.UpdateCounter(m.ID, *m.Delta)
		v, err := ms.GetCounter(m.ID)
		if err != nil {
			return err
		}
		m.Delta = &v.Value
		return nil
	default:
		return models.ErrInvalidMetricType
	}
}

func (ms *MemStorage) GetMetric(m *models.Metrics) error {
	if err := m.ValidateGet(); err != nil {
		return err
	}
	switch models.MetricType(m.MType) {
	case models.TypeGauge:
		v, err := ms.GetGauge(m.ID)
		if err != nil {
			return err
		}
		m.Value = &v.Value
		return nil
	case models.TypeCounter:
		v, err := ms.GetCounter(m.ID)
		if err != nil {
			return err
		}
		m.Delta = &v.Value
		return nil
	default:
		return models.ErrInvalidMetricType
	}
}
