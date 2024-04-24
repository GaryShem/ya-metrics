package memStorage

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/GaryShem/ya-metrics.git/internal/shared/storage"
	"github.com/GaryShem/ya-metrics.git/internal/shared/storage/metrics"
)

type MemStorageGaugeTestSuite struct {
	suite.Suite
	repo storage.Repository
}

func (s *MemStorageGaugeTestSuite) SetupSuite() {
	s.repo = NewMemStorage()
}

func TestMemStorageGaugeTestSuite(t *testing.T) {
	suite.Run(t, new(MemStorageGaugeTestSuite))
}

func (s *MemStorageGaugeTestSuite) TestGauge() {
	gauges := []*metrics.Gauge{
		metrics.NewGauge("foo", 3.14),
		metrics.NewGauge("foo", 4.14),
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
	gauges := map[string]*metrics.Gauge{
		"a": metrics.NewGauge("a", 3),
		"b": metrics.NewGauge("b", 1),
	}
	for _, g := range gauges {
		s.repo.UpdateGauge(g.Name, g.Value)
	}
	s.EqualValues(gauges, s.repo.GetGauges())
}
