package postgres

import (
	"github.com/GaryShem/ya-metrics.git/internal/shared/logging"
	"github.com/GaryShem/ya-metrics.git/internal/shared/storage/models"
)

func (s *SQLStorage) GetCounters() (map[string]models.Counter, error) {
	logging.Log.Infoln("Getting all SQL counters")
	result := make(map[string]models.Counter)
	queryTemplate := `SELECT * FROM counters`
	res, err := s.db.Query(queryTemplate)
	if err != nil {
		return nil, err
	}
	if res.Err() != nil {
		return nil, res.Err()
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

func (s *SQLStorage) GetCounter(metricName string) (*models.Counter, error) {
	logging.Log.Infoln("Getting SQL counter", metricName)
	queryTemplate := `SELECT * FROM counters WHERE id = $1`
	res := s.db.QueryRow(queryTemplate, metricName)
	counter := models.NewCounter("", 0)
	if err := res.Scan(&counter.Name, &counter.Value); err != nil {
		return nil, err
	}
	return counter, nil
}

func (s *SQLStorage) UpdateCounter(metricName string, delta int64) error {
	logging.Log.Infoln("Updating SQL counter", metricName, delta)
	queryTemplate := `INSERT INTO counters(id, val) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET val = counters.val + excluded.val`
	if _, err := s.db.Exec(queryTemplate, metricName, delta); err != nil {
		return err
	}
	return nil
}
