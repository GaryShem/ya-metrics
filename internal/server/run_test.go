package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/suite"

	"github.com/GaryShem/ya-metrics.git/internal/shared/storage"
)

type MetricHandlerSuite struct {
	suite.Suite
	server *httptest.Server
	repo   storage.Repository
}

func (s *MetricHandlerSuite) SetupSuite() {
	s.repo = storage.NewMemStorage()
	s.server = httptest.NewServer(MetricsRouter(s.repo))
}

func TestMetricHandlerSuite(t *testing.T) {
	suite.Run(t, new(MetricHandlerSuite))
}

func (s *MetricHandlerSuite) TestSetGauge() {
	reqlURL := "/update/gauge/foo/2"
	client := resty.New()
	url := s.server.URL + reqlURL
	response, err := client.R().Post(url)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, response.StatusCode())
}

func (s *MetricHandlerSuite) TestSetCounter() {
	reqURL := "/update/counter/foo/2"
	client := resty.New()
	url := s.server.URL + reqURL
	response, err := client.R().Post(url)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, response.StatusCode())
}

func (s *MetricHandlerSuite) TestGetGauge() {
	value := 3.14
	s.repo.UpdateGauge("foo", value)
	reqURL := "/value/gauge/foo"
	client := resty.New()
	url := s.server.URL + reqURL
	response, err := client.R().Get(url)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, response.StatusCode())
	s.Require().Equal(fmt.Sprint(value), string(response.Body()))
}

func (s *MetricHandlerSuite) TestGetCounter() {
	var value int64 = 42
	s.repo.UpdateCounter("foo", value)
	reqURL := "/value/counter/foo"
	client := resty.New()
	url := s.server.URL + reqURL
	response, err := client.R().Get(url)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, response.StatusCode())
	s.Require().Equal(fmt.Sprint(value), string(response.Body()))
}
