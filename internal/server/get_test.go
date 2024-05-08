package server

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"

	"github.com/GaryShem/ya-metrics.git/internal/shared/logging"
)

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

func (s *MetricHandlerSuite) TestGetIncorrectCounter() {
	reqURL := "/value/counter/nonexistent"
	client := resty.New()
	url := s.server.URL + reqURL
	response, err := client.R().Get(url)
	s.Require().NoError(err)
	logging.Log.Infoln(string(response.Body()))
	s.Require().Equal(http.StatusNotFound, response.StatusCode())
}
