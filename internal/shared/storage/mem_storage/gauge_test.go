package mem_storage

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/GaryShem/ya-metrics.git/internal/shared/storage/models"
)

type MemStorageGaugeTestSuite struct {
	suite.Suite
	repo models.Repository
}

func (s *MemStorageGaugeTestSuite) SetupSuite() {
	s.repo = NewMemStorage()
}

func TestMemStorageGaugeTestSuite(t *testing.T) {
	suite.Run(t, new(MemStorageGaugeTestSuite))
}

func (s *MemStorageGaugeTestSuite) TestGauge() {
	gauges := []*models.Gauge{
		models.NewGauge("foo", 3.14),
		models.NewGauge("foo", 4.14),
	}
	for _, gauge := range gauges {
		s.repo.UpdateGauge(gauge.Name, gauge.Value)
		rv, err := s.repo.GetGauge(gauge.Name)
		s.Require().NoError(err)
		s.InEpsilon(gauge.Value, rv.Value, 0.001)
	}

	_, err := s.repo.GetGauge("bar")
	s.Require().ErrorIs(err, ErrMetricNotFound)
}

func (s *MemStorageCounterTestSuite) TestGauges() {
	gauges := map[string]*models.Gauge{
		"a": models.NewGauge("a", 3),
		"b": models.NewGauge("b", 1),
	}
	for _, g := range gauges {
		s.repo.UpdateGauge(g.Name, g.Value)
	}
	s.EqualValues(gauges, s.repo.GetGauges())
}
