package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-resty/resty/v2"

	"github.com/GaryShem/ya-metrics.git/internal/shared/storage/models"
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

func (s *MetricHandlerSuite) TestSetCounterJSON() {
	reqURL := "/update/"
	client := resty.New()
	value := int64(42)
	m := models.Metrics{
		ID:    "foo",
		MType: string(models.TypeCounter),
		Delta: &value,
		Value: nil,
	}
	mJSON, err := json.Marshal(m)
	s.Require().NoError(err)
	url := s.server.URL + reqURL

	response, err := client.R().
		SetHeader("Content-Type", "application/json").SetBody(mJSON).Post(url)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, response.StatusCode())
	s.Equal(string(mJSON), string(response.Body()))
}

func (s *MetricHandlerSuite) TestSetGaugeJSON() {
	updateURL := "/update/"
	getURL := "/value/"
	client := resty.New()
	value := 3.14
	m := models.Metrics{
		ID:    "foo",
		MType: string(models.TypeGauge),
		Delta: nil,
		Value: &value,
	}
	mJSON, err := json.Marshal(m)
	s.Require().NoError(err)
	url := s.server.URL + updateURL

	response, err := client.R().
		SetHeader("Content-Type", "application/json").SetBody(mJSON).Post(url)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, response.StatusCode())
	s.Require().Equal("application/json", response.Header().Get("Content-Type"))
	s.Equal(string(mJSON), string(response.Body()))
	m.Value = nil

	url = s.server.URL + getURL
	response, err = client.R().
		SetHeader("Content-Type", "application/json").SetBody(mJSON).Post(url)
	s.Require().NoError(err)
	expectedStr := string(mJSON)
	gotStr := string(response.Body())
	s.Equal(expectedStr, gotStr)

}

func (s *MetricHandlerSuite) TestSetGaugeJSONIncorrectContentType() {
	reqURL := "/update/"
	client := resty.New()
	value := 3.14
	m := models.Metrics{
		ID:    "foo",
		MType: string(models.TypeGauge),
		Delta: nil,
		Value: &value,
	}
	mJSON, err := json.Marshal(m)
	s.Require().NoError(err)
	url := s.server.URL + reqURL

	response, err := client.R().
		SetHeader("Content-Type", "text/plain; charset=utf-8").SetBody(mJSON).Post(url)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusNotFound, response.StatusCode())
}
