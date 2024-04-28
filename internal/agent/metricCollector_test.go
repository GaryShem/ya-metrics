package agent

import (
	"encoding/json"
	"runtime"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/GaryShem/ya-metrics.git/internal/shared/storage/memorystorage"
	"github.com/GaryShem/ya-metrics.git/internal/shared/storage/models"
)

type MetricCollectorSuite struct {
	suite.Suite
	collector *MetricCollector
}

func (s *MetricCollectorSuite) SetupSuite() {
	s.collector = NewMetricCollector([]string{"Alloc"})
}

func TestMetricCollectorSuite(t *testing.T) {
	suite.Run(t, new(MetricCollectorSuite))
}

func (s *MetricCollectorSuite) TestMetricGetter() {
	runtimeMetrics := SupportedRuntimeMetrics()
	tests := make([]string, len(runtimeMetrics))
	copy(tests, runtimeMetrics)
	for _, tt := range tests {
		var rtm runtime.MemStats
		runtime.ReadMemStats(&rtm)
		got, err := Getter(&rtm, tt)
		s.Require().NoError(err)
		s.NotEqual(0, got)
	}
	invalidMetric := []string{
		"wololo",
	}
	for _, tt := range invalidMetric {
		var rtm runtime.MemStats
		runtime.ReadMemStats(&rtm)
		_, err := Getter(&rtm, tt)
		s.Require().Error(err)
	}
}

func (s *MetricCollectorSuite) TestCollectingMetrics() {
	collectedMetrics := []string{"Alloc"}
	expectedPollCount := int64(1)
	mc := NewMetricCollector([]string{"Alloc"})
	err := mc.CollectMetrics()
	s.Require().NoError(err)
	for _, name := range collectedMetrics {
		s.NotEqual(0, mc.Storage.GetGauges()[name])
	}

	pollCount, err := mc.Storage.GetCounter("PollCount")
	s.Require().NoError(err)
	s.Equal(expectedPollCount, pollCount.Value)
}

func (s *MetricCollectorSuite) TestIllegalMetric() {
	mc := &MetricCollector{
		Storage:                 memorystorage.NewMemStorage(),
		RuntimeGaugeMetricNames: []string{"Alloc"},
	}
	mc.Storage.UpdateGauge("Alloc", 2)
	mc.Storage.UpdateCounter("PollCount", 1)
	dump, err := mc.DumpMetrics()
	s.Require().NoError(err)
	dumpJSON, err := json.Marshal(dump)
	s.Require().NoError(err)
	allocValue := float64(2)
	pollCount := int64(1)
	want := []*models.Metrics{
		&models.Metrics{
			ID:    "Alloc",
			MType: string(models.TypeGauge),
			Delta: nil,
			Value: &allocValue,
		},
		&models.Metrics{
			ID:    "PollCount",
			MType: string(models.TypeCounter),
			Delta: &pollCount,
			Value: nil,
		},
	}
	wantJSON, err := json.Marshal(want)
	s.Equal(wantJSON, dumpJSON)
	s.Require().NoError(err)
	pc, err := mc.Storage.GetCounter("PollCount")
	s.Require().NoError(err)
	s.Equal(int64(0), pc.Value)
}
