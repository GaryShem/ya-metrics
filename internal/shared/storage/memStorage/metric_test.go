package memstorage

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/GaryShem/ya-metrics.git/internal/shared/storage"
	"github.com/GaryShem/ya-metrics.git/internal/shared/storage/metrics"
)

type MemStorageMetricsTestSuite struct {
	suite.Suite
	repo storage.Repository
}

func (s *MemStorageMetricsTestSuite) BeforeTest(suiteName, testName string) {
	s.repo = NewMemStorage()
}

func TestMemStorageTestSuite(t *testing.T) {
	suite.Run(t, new(MemStorageMetricsTestSuite))
}

func (s *MemStorageMetricsTestSuite) TestUpdateMetricInvalidMetricID() {
	m := &metrics.Metrics{
		ID:    "",
		MType: string(metrics.TypeCounter),
		Delta: nil,
		Value: nil,
	}
	err := s.repo.UpdateMetric(m)
	s.Require().ErrorIs(err, metrics.ErrInvalidMetricID)
}

func (s *MemStorageMetricsTestSuite) TestUpdateMetricInvalidMetricType() {
	m := &metrics.Metrics{
		ID:    "foo",
		MType: "bar",
		Delta: nil,
		Value: nil,
	}
	err := s.repo.UpdateMetric(m)
	s.Require().ErrorIs(err, metrics.ErrInvalidMetricType)
}

func (s *MemStorageMetricsTestSuite) TestUpdateMetricInvalidCounterValue() {
	m := &metrics.Metrics{
		ID:    "foo",
		MType: string(metrics.TypeCounter),
		Delta: nil,
		Value: nil,
	}
	err := s.repo.UpdateMetric(m)
	s.Require().ErrorIs(err, metrics.ErrInvalidMetricValue)
}

func (s *MemStorageMetricsTestSuite) TestUpdateMetricInvalidGaugeValue() {
	m := &metrics.Metrics{
		ID:    "foo",
		MType: string(metrics.TypeGauge),
		Delta: nil,
		Value: nil,
	}
	err := s.repo.UpdateMetric(m)
	s.Require().ErrorIs(err, metrics.ErrInvalidMetricValue)
}

func (s *MemStorageMetricsTestSuite) TestMetricValidGauge() {
	value := 3.14
	m := &metrics.Metrics{
		ID:    "foo",
		MType: string(metrics.TypeGauge),
		Delta: nil,
		Value: &value,
	}
	err := s.repo.UpdateMetric(m)
	s.Require().NoError(err)
	v, err := s.repo.GetGauge("foo")
	s.Require().NoError(err)
	s.InEpsilon(v.Value, value, 0.001)
	m2 := &metrics.Metrics{
		ID:    "foo",
		MType: string(metrics.TypeGauge),
		Delta: nil,
		Value: nil,
	}
	err = s.repo.GetMetric(m2)
	s.Require().NoError(err)
	s.EqualValues(m, m2)
}

func (s *MemStorageMetricsTestSuite) TestMetricValidCounter() {
	value := int64(42)
	m := &metrics.Metrics{
		ID:    "foo",
		MType: string(metrics.TypeCounter),
		Delta: &value,
		Value: nil,
	}
	err := s.repo.UpdateMetric(m)
	s.Require().NoError(err)
	v, err := s.repo.GetCounter("foo")
	s.Require().NoError(err)
	s.Equal(v.Value, value)

	m2 := &metrics.Metrics{
		ID:    "foo",
		MType: string(metrics.TypeCounter),
		Delta: nil,
		Value: nil,
	}
	err = s.repo.GetMetric(m2)
	s.Require().NoError(err)
	s.EqualValues(m, m2)
}
