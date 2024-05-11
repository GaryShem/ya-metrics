package postgres

import (
	"github.com/GaryShem/ya-metrics.git/internal/shared/logging"
	"github.com/GaryShem/ya-metrics.git/internal/shared/storage/models"
)

func (s *SqlStorage) GetCounters() (map[string]*models.Counter, error) {
	logging.Log.Infoln("Getting all SQL counters")
	result := make(map[string]*models.Counter)
	queryTemplate := `SELECT * FROM counters`
	res, err := s.db.Query(queryTemplate)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	for res.Next() {
		counter := models.NewCounter("", 0)
		err = res.Scan(&counter.Name, &counter.Value)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (s *SqlStorage) GetCounter(metricName string) (*models.Counter, error) {
	logging.Log.Infoln("Getting SQL counter", metricName)
	queryTemplate := `SELECT * FROM counters WHERE id = $1`
	res := s.db.QueryRow(queryTemplate, metricName)
	counter := models.NewCounter("", 0)
	if err := res.Scan(&counter.Name, &counter.Value); err != nil {
		return nil, err
	}
	return counter, nil
}

func (s *SqlStorage) UpdateCounter(metricName string, delta int64) error {
	logging.Log.Infoln("Updating SQL gauge", metricName, delta)
	queryTemplate := `INSERT INTO counters(id, val) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET val = excluded.val + $2`
	if _, err := s.db.Exec(queryTemplate, metricName, delta); err != nil {
		return err
	}
	return nil
}