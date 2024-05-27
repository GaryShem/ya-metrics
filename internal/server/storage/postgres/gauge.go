package postgres

import (
	"github.com/GaryShem/ya-metrics.git/internal/shared/logging"
	"github.com/GaryShem/ya-metrics.git/internal/shared/storage/models"
)

func (s *SQLStorage) GetGauges() (map[string]models.Gauge, error) {
	logging.Log.Infoln("Getting all SQL gauges")
	result := make(map[string]models.Gauge)
	queryTemplate := `SELECT * FROM gauges`
	res, err := s.db.Query(queryTemplate)
	if err != nil {
		return nil, err
	}
	if res.Err() != nil {
		return nil, res.Err()
	}
	defer res.Close()
	for res.Next() {
		gauge := models.NewGauge("", 0)
		err = res.Scan(&gauge.Name, &gauge.Value)
		if err != nil {
			return nil, err
		}
		result[gauge.Name] = *gauge
	}

	err = res.Close()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *SQLStorage) GetGauge(metricName string) (*models.Gauge, error) {
	logging.Log.Infoln("Getting SQL gauge", metricName)
	queryTemplate := `SELECT * FROM gauges WHERE id = $1`
	res := s.db.QueryRow(queryTemplate, metricName)
	gauge := models.NewGauge("", 0)
	if err := res.Scan(&gauge.Name, &gauge.Value); err != nil {
		return nil, err
	}
	return gauge, nil
}

func (s *SQLStorage) UpdateGauge(metricName string, value float64) error {
	logging.Log.Infoln("Updating SQL gauge", metricName, value)
	queryTemplate := `INSERT INTO gauges (id, val) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET val = $2`
	if _, err := s.db.Exec(queryTemplate, metricName, value); err != nil {
		return err
	}
	return nil
}
