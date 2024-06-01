package memorystorage

import (
	"github.com/GaryShem/ya-metrics.git/internal/shared/logging"
	"github.com/GaryShem/ya-metrics.git/internal/shared/storage/models"
)

func (ms *MemStorage) UpdateMetric(m *models.Metrics) error {
	if err := m.ValidateUpdate(); err != nil {
		return err
	}
	switch models.MetricType(m.MType) {
	case models.TypeGauge:
		err := ms.UpdateGauge(m.ID, *m.Value)
		if err != nil {
			return err
		}
		v, err := ms.GetGauge(m.ID)
		if err != nil {
			return err
		}
		m.Value = &v.Value
		return nil
	case models.TypeCounter:
		err := ms.UpdateCounter(m.ID, *m.Delta)
		if err != nil {
			return err
		}
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

func (ms *MemStorage) UpdateMetricBatch(metrics []models.Metrics) ([]models.Metrics, error) {
	metricNames := make([]string, 0)
	for _, m := range metrics {
		metricNames = append(metricNames, m.ID)
	}
	logging.Log.Infoln("updating metrics: ", metricNames)
	for _, m := range metrics {
		if err := ms.UpdateMetric(&m); err != nil {
			return nil, err
		}
	}
	result := make([]models.Metrics, 0)
	for _, m := range metrics {
		isDuplicate := false
		for _, r := range result {
			if r.MType == m.MType && r.ID == m.ID {
				isDuplicate = true
				break
			}
		}
		if isDuplicate {
			continue
		}
		err := ms.GetMetric(&m)
		if err != nil {
			return nil, err
		}
		result = append(result, m)
	}
	return result, nil
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

func (ms *MemStorage) ListMetrics() ([]models.Metrics, error) {
	result := make([]models.Metrics, 0)
	for _, v := range ms.GaugeMetrics {
		result = append(result, models.Metrics{
			ID:    v.Name,
			MType: v.Type,
			Delta: nil,
			Value: &v.Value,
		})
	}
	for _, v := range ms.CounterMetrics {
		result = append(result, models.Metrics{
			ID:    v.Name,
			MType: v.Type,
			Delta: &v.Value,
			Value: nil,
		})
	}
	return result, nil
}
