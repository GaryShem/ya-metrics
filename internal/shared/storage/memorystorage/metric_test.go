package memorystorage

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/GaryShem/ya-metrics.git/internal/shared/storage/models"
)

type MemStorageMetricsTestSuite struct {
	suite.Suite
	repo models.Repository
}

func (s *MemStorageMetricsTestSuite) BeforeTest(suiteName, testName string) {
	s.repo = NewMemStorage()
}

func TestMemStorageTestSuite(t *testing.T) {
	suite.Run(t, new(MemStorageMetricsTestSuite))
}

func (s *MemStorageMetricsTestSuite) TestUpdateMetricInvalidMetricID() {
	m := &models.Metrics{
		ID:    "",
		MType: string(models.TypeCounter),
		Delta: nil,
		Value: nil,
	}
	err := s.repo.UpdateMetric(m)
	s.Require().ErrorIs(err, models.ErrInvalidMetricID)
}

func (s *MemStorageMetricsTestSuite) TestUpdateMetricInvalidMetricType() {
	m := &models.Metrics{
		ID:    "foo",
		MType: "bar",
		Delta: nil,
		Value: nil,
	}
	err := s.repo.UpdateMetric(m)
	s.Require().ErrorIs(err, models.ErrInvalidMetricType)
}

func (s *MemStorageMetricsTestSuite) TestUpdateMetricInvalidCounterValue() {
	m := &models.Metrics{
		ID:    "foo",
		MType: string(models.TypeCounter),
		Delta: nil,
		Value: nil,
	}
	err := s.repo.UpdateMetric(m)
	s.Require().ErrorIs(err, models.ErrInvalidMetricValue)
}

func (s *MemStorageMetricsTestSuite) TestUpdateMetricInvalidGaugeValue() {
	m := &models.Metrics{
		ID:    "foo",
		MType: string(models.TypeGauge),
		Delta: nil,
		Value: nil,
	}
	err := s.repo.UpdateMetric(m)
	s.Require().ErrorIs(err, models.ErrInvalidMetricValue)
}

func (s *MemStorageMetricsTestSuite) TestMetricValidGauge() {
	value := 3.14
	m := &models.Metrics{
		ID:    "foo",
		MType: string(models.TypeGauge),
		Delta: nil,
		Value: &value,
	}
	err := s.repo.UpdateMetric(m)
	s.Require().NoError(err)
	v, err := s.repo.GetGauge("foo")
	s.Require().NoError(err)
	s.InEpsilon(v.Value, value, 0.001)
	m2 := &models.Metrics{
		ID:    "foo",
		MType: string(models.TypeGauge),
		Delta: nil,
		Value: nil,
	}
	err = s.repo.GetMetric(m2)
	s.Require().NoError(err)
	s.EqualValues(m, m2)
}

func (s *MemStorageMetricsTestSuite) TestMetricValidCounter() {
	value := int64(42)
	m := &models.Metrics{
		ID:    "foo",
		MType: string(models.TypeCounter),
		Delta: &value,
		Value: nil,
	}
	err := s.repo.UpdateMetric(m)
	s.Require().NoError(err)
	v, err := s.repo.GetCounter("foo")
	s.Require().NoError(err)
	s.Equal(v.Value, value)

	m2 := &models.Metrics{
		ID:    "foo",
		MType: string(models.TypeCounter),
		Delta: nil,
		Value: nil,
	}
	err = s.repo.GetMetric(m2)
	s.Require().NoError(err)
	s.EqualValues(m, m2)
}
