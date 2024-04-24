package memstorage

import "github.com/GaryShem/ya-metrics.git/internal/shared/storage/metrics"

func (ms *MemStorage) UpdateMetric(m *metrics.Metrics) error {
	if err := m.ValidateUpdate(); err != nil {
		return err
	}
	switch metrics.MetricType(m.MType) {
	case metrics.TypeGauge:
		ms.UpdateGauge(m.ID, *m.Value)
		v, err := ms.GetGauge(m.ID)
		if err != nil {
			return err
		}
		m.Value = &v.Value
		return nil
	case metrics.TypeCounter:
		ms.UpdateCounter(m.ID, *m.Delta)
		v, err := ms.GetCounter(m.ID)
		if err != nil {
			return err
		}
		m.Delta = &v.Value
		return nil
	default:
		return metrics.ErrInvalidMetricType
	}
}

func (ms *MemStorage) GetMetric(m *metrics.Metrics) error {
	if err := m.ValidateGet(); err != nil {
		return err
	}
	switch metrics.MetricType(m.MType) {
	case metrics.TypeGauge:
		v, err := ms.GetGauge(m.ID)
		if err != nil {
			return err
		}
		m.Value = &v.Value
		return nil
	case metrics.TypeCounter:
		v, err := ms.GetCounter(m.ID)
		if err != nil {
			return err
		}
		m.Delta = &v.Value
		return nil
	default:
		return metrics.ErrInvalidMetricType
	}
}
