package postgres

import (
	"fmt"
	"strings"

	"github.com/GaryShem/ya-metrics.git/internal/shared/logging"
	"github.com/GaryShem/ya-metrics.git/internal/shared/storage/models"
)

func (s *SQLStorage) UpdateMetric(m *models.Metrics) error {
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

func (s *SQLStorage) UpdateMetricBatch(metrics []*models.Metrics) ([]*models.Metrics, error) {
	result := make([]*models.Metrics, 0)
	if len(metrics) == 0 {
		return result, nil
	}
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	gauges := make([]*models.Metrics, 0)
	counters := make([]*models.Metrics, 0)
	for _, m := range metrics {
		switch m.MType {
		case string(models.TypeGauge):
			isDuplicate := false
			for _, g := range gauges {
				if g.ID == m.ID {
					isDuplicate = true
					g.Value = m.Value
					break
				}
			}
			if !isDuplicate {
				gauges = append(gauges, m)
			}
		case string(models.TypeCounter):
			isDuplicate := false
			for _, c := range counters {
				if c.ID == m.ID {
					isDuplicate = true
					value := *c.Delta + *m.Delta
					c.Delta = &value
					break
				}
			}
			if !isDuplicate {
				counters = append(counters, m)
			}
		default:
			return nil, models.ErrInvalidMetricType
		}
	}
	gauges, err = s.updateGauges(gauges)
	if err != nil {
		return nil, fmt.Errorf("failed to update gauges: %w", err)
	}
	counters, err = s.updateCounters(counters)
	if err != nil {
		return nil, fmt.Errorf("failed to update counters: %w", err)
	}
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	result = make([]*models.Metrics, len(gauges)+len(counters))
	result = append(result, gauges...)
	result = append(result, counters...)
	return result, nil
}

func (s *SQLStorage) updateGauges(metrics []*models.Metrics) ([]*models.Metrics, error) {
	result := make([]*models.Metrics, 0)
	if len(metrics) == 0 {
		return result, nil
	}
	templatePieces := make([]string, 0)
	valuePieces := make([]interface{}, 0)
	i := 1
	for _, m := range metrics {
		templatePieces = append(templatePieces, fmt.Sprintf("($%d, $%d)", i, i+1))
		valuePieces = append(valuePieces, m.ID, *m.Value)
		i += 2
	}
	valuesString := strings.Join(templatePieces, ", ")

	queryTemplate := fmt.Sprintf(`INSERT INTO gauges(id, val) VALUES %s ON CONFLICT (id) DO UPDATE SET val = excluded.val RETURNING *`, valuesString)
	logging.Log.Infoln("gauge update sql:", queryTemplate, valuePieces)
	rows, err := s.db.Query(queryTemplate, valuePieces...)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, err
	}
	for rows.Next() {
		row := models.NewGauge("", 0)
		if err = rows.Scan(&row.Name, &row.Value); err != nil {
			return nil, err
		}
		result = append(result, &models.Metrics{
			MType: string(models.TypeGauge),
			ID:    row.Name,
			Delta: nil,
			Value: &row.Value,
		})
	}

	return result, nil
}

func (s *SQLStorage) updateCounters(metrics []*models.Metrics) ([]*models.Metrics, error) {
	result := make([]*models.Metrics, 0)
	if len(metrics) == 0 {
		return result, nil
	}
	templatePieces := make([]string, 0)
	valuePieces := make([]interface{}, 0)
	i := 1
	for _, m := range metrics {
		templatePieces = append(templatePieces, fmt.Sprintf("($%d, $%d)", i, i+1))
		valuePieces = append(valuePieces, m.ID, *m.Delta)
		i += 2
	}
	valuesString := strings.Join(templatePieces, ", ")

	queryTemplate := fmt.Sprintf(`INSERT INTO counters(id, val) VALUES %s ON CONFLICT (id) DO UPDATE SET val = counters.val + excluded.val RETURNING *`, valuesString)
	logging.Log.Infoln("counter update sql:", queryTemplate, valuePieces)
	rows, err := s.db.Query(queryTemplate, valuePieces...)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, err
	}
	for rows.Next() {
		row := models.NewCounter("", 0)
		if err = rows.Scan(&row.Name, &row.Value); err != nil {
			return nil, err
		}
		result = append(result, &models.Metrics{
			MType: string(models.TypeCounter),
			ID:    row.Name,
			Delta: &row.Value,
			Value: nil,
		})
	}

	return result, nil
}

func (s *SQLStorage) GetMetric(m *models.Metrics) error {
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

func (s *SQLStorage) ListMetrics() ([]*models.Metrics, error) {
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
