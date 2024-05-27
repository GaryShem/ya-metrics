package memorystorage

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/GaryShem/ya-metrics.git/internal/server/storage/repository"
	"github.com/GaryShem/ya-metrics.git/internal/shared/storage/models"
)

type MemStorageCounterTestSuite struct {
	suite.Suite
	repo repository.Repository
}

func (s *MemStorageCounterTestSuite) BeforeTest(_, _ string) {
	s.repo = NewMemStorage()
}

func TestMemStorageCounterTestSuite(t *testing.T) {
	suite.Run(t, new(MemStorageCounterTestSuite))
}

func (s *MemStorageCounterTestSuite) TestCounter() {
	counters := []*models.Counter{
		models.NewCounter("foo", 1),
		models.NewCounter("foo", 1),
	}
	sum := int64(0)
	for _, c := range counters {
		sum += c.Value
		err := s.repo.UpdateCounter(c.Name, c.Value)
		s.Require().NoError(err)
		rv, err := s.repo.GetCounter(c.Name)
		s.Require().NoError(err)
		s.Equal(sum, rv.Value)
	}

	_, err := s.repo.GetCounter("bar")
	s.Require().ErrorIs(err, repository.ErrMetricNotFound)
}

func (s *MemStorageCounterTestSuite) TestCounters() {
	counters := map[string]models.Counter{
		"a": *models.NewCounter("a", 1),
		"b": *models.NewCounter("b", 1),
	}
	for _, c := range counters {
		err := s.repo.UpdateCounter(c.Name, c.Value)
		s.Require().NoError(err)
	}
	value, err := s.repo.GetCounters()
	s.Require().NoError(err)
	s.EqualValues(counters, value)
}
