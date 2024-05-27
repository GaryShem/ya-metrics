package postgres

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/GaryShem/ya-metrics.git/internal/server/config"
	"github.com/GaryShem/ya-metrics.git/internal/server/storage/repository"
	"github.com/GaryShem/ya-metrics.git/internal/shared/logging"
	"github.com/GaryShem/ya-metrics.git/internal/shared/storage/models"
)

type SQLStorageMetricsTestSuite struct {
	suite.Suite
	repo       repository.Repository
	sqlStorage *SQLStorage
	connString string
}

func (s *SQLStorageMetricsTestSuite) SetupTest() {
	err := logging.InitializeZapLogger("Info")
	if err != nil {
		panic(err)
	}
	repo, err := NewSQLStorage(s.connString, true)
	if err != nil {
		panic(err)
	}
	s.repo = repo
	s.sqlStorage = repo
}

func (s *SQLStorageMetricsTestSuite) BeforeTest(_, _ string) {
	err := s.sqlStorage.Reset()
	if err != nil {
		panic(err)
	}
}

func TestMemStorageTestSuite(t *testing.T) {
	sf := config.ServerFlags{}
	config.ParseFlags(&sf)
	if sf.DBString == "" {
		t.Skip("No db string provided")
	}
	suite.Run(t, &SQLStorageMetricsTestSuite{
		connString: sf.DBString,
	})
}

func (s *SQLStorageMetricsTestSuite) TestUpdateMetricInvalidMetricID() {
	m := &models.Metrics{
		ID:    "",
		MType: string(models.TypeCounter),
		Delta: nil,
		Value: nil,
	}
	err := s.repo.UpdateMetric(m)
	s.Require().ErrorIs(err, models.ErrInvalidMetricID)
}

func (s *SQLStorageMetricsTestSuite) TestUpdateMetricInvalidMetricType() {
	m := &models.Metrics{
		ID:    "foo",
		MType: "bar",
		Delta: nil,
		Value: nil,
	}
	err := s.repo.UpdateMetric(m)
	s.Require().ErrorIs(err, models.ErrInvalidMetricType)
}

func (s *SQLStorageMetricsTestSuite) TestUpdateMetricInvalidCounterValue() {
	m := &models.Metrics{
		ID:    "foo",
		MType: string(models.TypeCounter),
		Delta: nil,
		Value: nil,
	}
	err := s.repo.UpdateMetric(m)
	s.Require().ErrorIs(err, models.ErrInvalidMetricValue)
}

func (s *SQLStorageMetricsTestSuite) TestUpdateMetricInvalidGaugeValue() {
	m := &models.Metrics{
		ID:    "foo",
		MType: string(models.TypeGauge),
		Delta: nil,
		Value: nil,
	}
	err := s.repo.UpdateMetric(m)
	s.Require().ErrorIs(err, models.ErrInvalidMetricValue)
}

func (s *SQLStorageMetricsTestSuite) TestMetricValidGauge() {
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

func (s *SQLStorageMetricsTestSuite) TestMetricValidCounter() {
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

func (s *SQLStorageMetricsTestSuite) TestDuplicateGaugeBatch() {
	value1 := 3.14
	value2 := 42.14
	m := []models.Metrics{
		{
			ID:    "foo",
			MType: string(models.TypeGauge),
			Delta: nil,
			Value: &value1,
		},
		{
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

func (s *SQLStorageMetricsTestSuite) TestDuplicateCounterBatch() {
	delta1 := int64(42)
	deltaDouble := int64(100)
	delta2 := int64(142)
	m := []models.Metrics{
		{
			ID:    "foo",
			MType: string(models.TypeCounter),
			Delta: &delta1,
			Value: nil,
		},
		{
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
