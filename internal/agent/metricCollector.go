package agent

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"

	"github.com/GaryShem/ya-metrics.git/internal/shared/storage/models"
)

func SupportedRuntimeMetrics() []string {
	return []string{
		"Alloc",
		"BuckHashSys",
		"Frees",
		"GCCPUFraction",
		"GCSys",
		"HeapAlloc",
		"HeapIdle",
		"HeapInuse",
		"HeapObjects",
		"HeapReleased",
		"HeapSys",
		"LastGC",
		"Lookups",
		"MCacheInuse",
		"MCacheSys",
		"MSpanInuse",
		"MSpanSys",
		"Mallocs",
		"NextGC",
		"NumForcedGC",
		"NumGC",
		"OtherSys",
		"PauseTotalNs",
		"StackInuse",
		"StackSys",
		"Sys",
		"TotalAlloc",
	}
}

//nolint:funlen
func Getter(m *runtime.MemStats, metricName string) (float64, error) {
	switch metricName {
	case "Alloc":
		return float64(m.Alloc), nil
	case "BuckHashSys":
		return float64(m.BuckHashSys), nil
	case "Frees":
		return float64(m.Frees), nil
	case "GCCPUFraction":
		return m.GCCPUFraction, nil
	case "GCSys":
		return float64(m.GCSys), nil
	case "HeapAlloc":
		return float64(m.HeapAlloc), nil
	case "HeapIdle":
		return float64(m.HeapIdle), nil
	case "HeapInuse":
		return float64(m.HeapInuse), nil
	case "HeapObjects":
		return float64(m.HeapObjects), nil
	case "HeapReleased":
		return float64(m.HeapReleased), nil
	case "HeapSys":
		return float64(m.HeapSys), nil
	case "LastGC":
		return float64(m.LastGC), nil
	case "Lookups":
		return float64(m.Lookups), nil
	case "MCacheInuse":
		return float64(m.MCacheInuse), nil
	case "MCacheSys":
		return float64(m.MCacheSys), nil
	case "MSpanInuse":
		return float64(m.MSpanInuse), nil
	case "MSpanSys":
		return float64(m.MSpanSys), nil
	case "Mallocs":
		return float64(m.Mallocs), nil
	case "NextGC":
		return float64(m.NextGC), nil
	case "NumForcedGC":
		return float64(m.NumForcedGC), nil
	case "NumGC":
		return float64(m.NumGC), nil
	case "OtherSys":
		return float64(m.OtherSys), nil
	case "PauseTotalNs":
		return float64(m.PauseTotalNs), nil
	case "StackInuse":
		return float64(m.StackInuse), nil
	case "StackSys":
		return float64(m.StackSys), nil
	case "Sys":
		return float64(m.Sys), nil
	case "TotalAlloc":
		return float64(m.TotalAlloc), nil
	}
	return 0, fmt.Errorf("metric not supported: %s", metricName)
}

type MetricCollector struct {
	Gauges                  map[string]float64
	Counters                map[string]int64
	RuntimeGaugeMetricNames []string
	mu                      sync.Mutex
}

func NewMetricCollector(gaugeMetrics []string) *MetricCollector {
	tmpGaugeMetrics := make([]string, len(gaugeMetrics))
	copy(tmpGaugeMetrics, gaugeMetrics)
	return &MetricCollector{
		Gauges:                  make(map[string]float64, len(tmpGaugeMetrics)),
		Counters:                make(map[string]int64, 1),
		RuntimeGaugeMetricNames: tmpGaugeMetrics,
	}
}

func (m *MetricCollector) CollectMetrics() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)
	for _, gaugeMetric := range m.RuntimeGaugeMetricNames {
		value, err := Getter(&rtm, gaugeMetric)
		if err != nil {
			return fmt.Errorf("could not read metric %s: %w", gaugeMetric, err)
		}
		m.Gauges[gaugeMetric] = value
	}
	m.Counters["PollCount"] += int64(len(m.RuntimeGaugeMetricNames))
	m.Gauges["RandomValue"] = rand.Float64() //nolint:gosec // do not need cryptography on this
	return nil
}

func (m *MetricCollector) DumpMetrics() ([]*models.Metrics, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	result := make([]*models.Metrics, 0)
	for id, value := range m.Gauges {
		metricValue := value
		result = append(result, &models.Metrics{
			ID:    id,
			MType: string(models.TypeGauge),
			Delta: nil,
			Value: &metricValue,
		})
	}
	for id, value := range m.Counters {
		metricValue := value
		result = append(result, &models.Metrics{
			ID:    id,
			MType: string(models.TypeCounter),
			Delta: &metricValue,
			Value: nil,
		})
	}

	m.Counters["PollCount"] = 0

	return result, nil
}
