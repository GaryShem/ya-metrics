package metrics

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type MetricsTestSuite struct {
	suite.Suite
}

func TestMetricsTestSuite(t *testing.T) {
	suite.Run(t, new(MetricsTestSuite))
}

func (s *MetricsTestSuite) TestValidateMetricInvalidID() {
	metric := &Metrics{
		ID:    "",
		MType: string(TypeCounter),
		Delta: nil,
		Value: nil,
	}
	err := metric.ValidateUpdate()
	s.Require().ErrorIs(err, ErrInvalidMetricID)
	err = metric.ValidateGet()
	s.Require().ErrorIs(err, ErrInvalidMetricID)
}

func (s *MetricsTestSuite) TestValidateMetricInvalidMType() {
	metric := &Metrics{
		ID:    "foo",
		MType: "",
		Delta: nil,
		Value: nil,
	}
	err := metric.ValidateUpdate()
	s.Require().ErrorIs(err, ErrInvalidMetricType)
	err = metric.ValidateGet()
	s.Require().ErrorIs(err, ErrInvalidMetricType)
}

func (s *MetricsTestSuite) TestValidateMetricInvalidDelta() {
	metric := &Metrics{
		ID:    "foo",
		MType: string(TypeCounter),
		Delta: nil,
		Value: nil,
	}
	err := metric.ValidateUpdate()
	s.Require().ErrorIs(err, ErrInvalidMetricValue)
	err = metric.ValidateGet()
	s.Require().NoError(err)
}

func (s *MetricsTestSuite) TestValidateMetricInvalidValue() {
	metric := &Metrics{
		ID:    "foo",
		MType: string(TypeGauge),
		Delta: nil,
		Value: nil,
	}
	err := metric.ValidateUpdate()
	s.Require().ErrorIs(err, ErrInvalidMetricValue)
	err = metric.ValidateGet()
	s.Require().NoError(err)
}
