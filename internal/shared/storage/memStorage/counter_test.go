package memstorage

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/GaryShem/ya-metrics.git/internal/shared/storage"
	"github.com/GaryShem/ya-metrics.git/internal/shared/storage/metrics"
)

type MemStorageCounterTestSuite struct {
	suite.Suite
	repo storage.Repository
}

func (s *MemStorageCounterTestSuite) BeforeTest(_, _ string) {
	s.repo = NewMemStorage()
}

func TestMemStorageCounterTestSuite(t *testing.T) {
	suite.Run(t, new(MemStorageCounterTestSuite))
}

func (s *MemStorageCounterTestSuite) TestCounter() {
	counters := []*metrics.Counter{
		metrics.NewCounter("foo", 1),
		metrics.NewCounter("foo", 1),
	}
	sum := int64(0)
	for _, c := range counters {
		sum += c.Value
		s.repo.UpdateCounter(c.Name, c.Value)
		rv, err := s.repo.GetCounter(c.Name)
		s.Require().NoError(err)
		s.Equal(sum, rv.Value)
	}
	err := s.repo.ResetCounter("foo")
	s.Require().NoError(err)
	m, err := s.repo.GetCounter("foo")
	s.Require().NoError(err)
	s.Equal(m.Value, int64(0))

	_, err = s.repo.GetCounter("bar")
	s.Require().ErrorIs(err, ErrMetricNotFound)
}

func (s *MemStorageCounterTestSuite) TestCounters() {
	counters := map[string]*metrics.Counter{
		"a": metrics.NewCounter("a", 1),
		"b": metrics.NewCounter("b", 1),
	}
	for _, c := range counters {
		s.repo.UpdateCounter(c.Name, c.Value)
	}
	s.EqualValues(counters, s.repo.GetCounters())
}
