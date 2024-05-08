package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"

	"github.com/GaryShem/ya-metrics.git/internal/shared/logging"
	"github.com/GaryShem/ya-metrics.git/internal/shared/storage/models"
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

func (s *MetricHandlerSuite) TestGetIncorrectGauge() {
	reqURL := "/value/gauge/nonexistent"
	client := resty.New()
	url := s.server.URL + reqURL
	response, err := client.R().Get(url)
	s.Require().NoError(err)
	logging.Log.Infoln(string(response.Body()))
	s.Require().Equal(http.StatusNotFound, response.StatusCode())
}
func (s *MetricHandlerSuite) TestGetGaugeEmptyName() {
	reqURL := "/value/gauge/"
	client := resty.New()
	url := s.server.URL + reqURL
	response, err := client.R().Get(url)
	s.Require().NoError(err)
	logging.Log.Infoln(string(response.Body()))
	s.Require().Equal(http.StatusNotFound, response.StatusCode())
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

func (s *MetricHandlerSuite) TestGetCounters() {
	var value int64 = 42
	s.repo.UpdateCounter("foo", value)
	reqURL := "/"
	client := resty.New()
	url := s.server.URL + reqURL
	response, err := client.R().Get(url)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, response.StatusCode())
	rs := struct {
		GaugeMetrics   map[string]*models.Gauge   `json:"gaugeMetrics"`
		CounterMetrics map[string]*models.Counter `json:"counterMetrics"`
	}{}
	responseString := string(response.Body())
	logging.Log.Infoln(responseString)
	err = json.Unmarshal([]byte(responseString), &rs)
	s.Require().NoError(err)
	s.Assert().Equal(len(rs.GaugeMetrics), 0)
	s.Assert().Equal(len(rs.CounterMetrics), 1)
	metric, ok := rs.CounterMetrics["foo"]
	s.Require().True(ok)
	s.Assert().Equal(value, metric.Value)
}

func (s *MetricHandlerSuite) TestGetCounterMetricJSON() {
	value := int64(42)
	s.repo.UpdateCounter("foo", value)

	m := models.Metrics{
		ID:    "foo",
		MType: string(models.TypeCounter),
		Delta: nil,
		Value: nil,
	}
	mJSON, err := json.Marshal(m)
	s.Require().NoError(err)

	reqURL := "/value"
	client := resty.New()
	url := s.server.URL + reqURL
	response, err := client.R().SetHeader("Content-Type", "application/json").SetBody(mJSON).Post(url)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, response.StatusCode())
	responseMetric := models.Metrics{}
	err = json.Unmarshal(response.Body(), &responseMetric)
	s.Require().NoError(err)
	s.Assert().EqualValues(value, *responseMetric.Delta)
	s.Equal(responseMetric.ID, m.ID)
	s.Equal(responseMetric.MType, m.MType)
}

func (s *MetricHandlerSuite) TestGetCounterMetricJSONInvalid() {
	m := models.Metrics{
		ID:    "foo",
		MType: string(models.TypeCounter),
		Delta: nil,
		Value: nil,
	}
	mJSON, err := json.Marshal(m)
	s.Require().NoError(err)

	reqURL := "/value"
	client := resty.New()
	url := s.server.URL + reqURL
	response, err := client.R().SetHeader("Content-Type", "application/json").SetBody(mJSON).Post(url)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusNotFound, response.StatusCode())
}

func (s *MetricHandlerSuite) TestGetGaugeMetricJSON() {
	value := 3.14
	s.repo.UpdateGauge("foo", value)

	m := models.Metrics{
		ID:    "foo",
		MType: string(models.TypeGauge),
		Delta: nil,
		Value: nil,
	}
	mJSON, err := json.Marshal(m)
	s.Require().NoError(err)

	reqURL := "/value"
	client := resty.New()
	url := s.server.URL + reqURL
	response, err := client.R().SetHeader("Content-Type", "application/json").SetBody(mJSON).Post(url)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, response.StatusCode())
	responseMetric := models.Metrics{}
	err = json.Unmarshal(response.Body(), &responseMetric)
	s.Require().NoError(err)
	s.Assert().EqualValues(value, *responseMetric.Value)
	s.Equal(responseMetric.ID, m.ID)
	s.Equal(responseMetric.MType, m.MType)
}

func (s *MetricHandlerSuite) TestGetGaugeMetricJSONInvalid() {
	m := models.Metrics{
		ID:    "foo",
		MType: string(models.TypeGauge),
		Delta: nil,
		Value: nil,
	}
	mJSON, err := json.Marshal(m)
	s.Require().NoError(err)

	reqURL := "/value"
	client := resty.New()
	url := s.server.URL + reqURL
	response, err := client.R().SetHeader("Content-Type", "application/json").SetBody(mJSON).Post(url)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusNotFound, response.StatusCode())
}
