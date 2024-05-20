package postgres

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/GaryShem/ya-metrics.git/internal/server/storage/repository"
	"github.com/GaryShem/ya-metrics.git/internal/shared/logging"
	"github.com/GaryShem/ya-metrics.git/internal/shared/storage/models"
)

type SqlStorageMetricsTestSuite struct {
	suite.Suite
	repo repository.Repository
}

func (s *SqlStorageMetricsTestSuite) SetupTest() {
	err := logging.InitializeZapLogger("Info")
	if err != nil {
		panic(err)
	}
}

func (s *SqlStorageMetricsTestSuite) BeforeTest(suiteName, testName string) {
	repo, err := NewSQLStorage("host=localhost user=postgres password=1231 dbname=postgres sslmode=disable", true)
	if err != nil {
		panic(err)
	}
	s.repo = repo
}

func TestMemStorageTestSuite(t *testing.T) {
	suite.Run(t, new(SqlStorageMetricsTestSuite))
}

func (s *SqlStorageMetricsTestSuite) TestUpdateMetricInvalidMetricID() {
	m := &models.Metrics{
		ID:    "",
		MType: string(models.TypeCounter),
		Delta: nil,
		Value: nil,
	}
	err := s.repo.UpdateMetric(m)
	s.Require().ErrorIs(err, models.ErrInvalidMetricID)
}

func (s *SqlStorageMetricsTestSuite) TestUpdateMetricInvalidMetricType() {
	m := &models.Metrics{
		ID:    "foo",
		MType: "bar",
		Delta: nil,
		Value: nil,
	}
	err := s.repo.UpdateMetric(m)
	s.Require().ErrorIs(err, models.ErrInvalidMetricType)
}

func (s *SqlStorageMetricsTestSuite) TestUpdateMetricInvalidCounterValue() {
	m := &models.Metrics{
		ID:    "foo",
		MType: string(models.TypeCounter),
		Delta: nil,
		Value: nil,
	}
	err := s.repo.UpdateMetric(m)
	s.Require().ErrorIs(err, models.ErrInvalidMetricValue)
}

func (s *SqlStorageMetricsTestSuite) TestUpdateMetricInvalidGaugeValue() {
	m := &models.Metrics{
		ID:    "foo",
		MType: string(models.TypeGauge),
		Delta: nil,
		Value: nil,
	}
	err := s.repo.UpdateMetric(m)
	s.Require().ErrorIs(err, models.ErrInvalidMetricValue)
}

func (s *SqlStorageMetricsTestSuite) TestMetricValidGauge() {
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

func (s *SqlStorageMetricsTestSuite) TestMetricValidCounter() {
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

func (s *SqlStorageMetricsTestSuite) TestDuplicateGaugeBatch() {
	value1 := 3.14
	value2 := 42.14
	m := []models.Metrics{
		models.Metrics{
			ID:    "foo",
			MType: string(models.TypeGauge),
			Delta: nil,
			Value: &value1,
		},
		models.Metrics{
			ID:    "foo",
			MType: string(models.TypeGauge),
			Delta: nil,
			Value: &value2,
		},
	}
	canonResult := []models.Metrics{m[1]}

	metrics, err := s.repo.UpdateMetricBatch(m)
	s.Require().NoError(err)
	s.Equal(len(canonResult), len(metrics))
	s.Equal(canonResult[0], metrics[0])
	v, err := s.repo.GetGauge("foo")
	s.Require().NoError(err)
	s.InEpsilon(v.Value, value2, 0.001)
	m2 := models.Metrics{
		ID:    "foo",
		MType: string(models.TypeGauge),
		Delta: nil,
		Value: nil,
	}
	err = s.repo.GetMetric(&m2)
	s.Require().NoError(err)
	s.EqualValues(m[len(m)-1], m2)
}

func (s *SqlStorageMetricsTestSuite) TestDuplicateCounterBatch() {
	delta1 := int64(42)
	deltaDouble := int64(100)
	delta2 := int64(142)
	m := []models.Metrics{
		models.Metrics{
			ID:    "foo",
			MType: string(models.TypeCounter),
			Delta: &delta1,
			Value: nil,
		},
		models.Metrics{
			ID:    "foo",
			MType: string(models.TypeCounter),
			Delta: &deltaDouble,
			Value: nil,
		},
	}
	canonResult := models.Metrics{
		ID:    "foo",
		MType: string(models.TypeCounter),
		Delta: &delta2,
		Value: nil,
	}

	metrics, err := s.repo.UpdateMetricBatch(m)
	s.Require().NoError(err)
	s.Equal(1, len(metrics))
	s.Equal(canonResult, metrics[0])
	v, err := s.repo.GetCounter("foo")
	s.Require().NoError(err)
	s.Equal(v.Value, delta2)
	m2 := models.Metrics{
		ID:    "foo",
		MType: string(models.TypeCounter),
		Delta: nil,
		Value: nil,
	}
	err = s.repo.GetMetric(&m2)
	s.Require().NoError(err)
	s.EqualValues(canonResult, m2)
}
