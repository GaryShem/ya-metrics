package postgres

import (
	"github.com/GaryShem/ya-metrics.git/internal/shared/storage/models"
)

func (s *SqlStorage) UpdateMetric(m *models.Metrics) error {
	switch m.MType {
	case string(models.TypeGauge):
		if err := s.UpdateGauge(m.ID, *m.Value); err != nil {
			return err
		}
	case string(models.TypeCounter):
		if err := s.UpdateCounter(m.ID, *m.Delta); err != nil {
			return err
		}
	default:
		return models.ErrInvalidMetricType
	}
	return nil
}

func (s *SqlStorage) GetMetric(m *models.Metrics) error {
	switch m.MType {
	case string(models.TypeGauge):
		gauge, err := s.GetGauge(m.ID)
		if err != nil {
			return err
		}
		m.Value = &gauge.Value
	case string(models.TypeCounter):
		counter, err := s.GetCounter(m.ID)
		if err != nil {
			return err
		}
		m.Delta = &counter.Value
	default:
		return models.ErrInvalidMetricType
	}
	return nil
}

func (s *SqlStorage) ListMetrics() ([]*models.Metrics, error) {
	result := make([]*models.Metrics, 0)
	gauges, err := s.GetGauges()
	if err != nil {
		return nil, err
	}
	for _, g := range gauges {
		result = append(result, &models.Metrics{
			ID:    g.Name,
			MType: string(models.TypeGauge),
			Delta: nil,
			Value: &g.Value,
		})
	}
	counters, err := s.GetCounters()
	if err != nil {
		return nil, err
	}
	for _, c := range counters {
		result = append(result, &models.Metrics{
			ID:    c.Name,
			MType: string(models.TypeCounter),
			Delta: &c.Value,
			Value: nil,
		})
	}
	return result, nil
}
