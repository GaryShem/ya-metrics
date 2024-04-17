package server

import (
	"net/http"

	"github.com/go-resty/resty/v2"
)

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

func (s *MetricHandlerSuite) TestSetGaugeGet() {
	reqlURL := "/update/gauge/foo/2"
	client := resty.New()
	url := s.server.URL + reqlURL
	response, err := client.R().Get(url)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusBadRequest, response.StatusCode())
}

func (s *MetricHandlerSuite) TestSetCounterGet() {
	reqURL := "/update/counter/foo/2"
	client := resty.New()
	url := s.server.URL + reqURL
	response, err := client.R().Get(url)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusBadRequest, response.StatusCode())
}
