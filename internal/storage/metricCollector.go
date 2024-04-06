package storage

import (
	"fmt"
	"math/rand"
	"runtime"
)

var RuntimeMetrics []string = []string{
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
	"OtherSys",
	"PauseTotalNs",
	"StackInuse",
	"StackSys",
	"Sys",
	"TotalAlloc",
}

func Getter(m *runtime.MemStats, metricName string) (float64, error) {
	switch metricName {
	case "Alloc":
		return float64(m.Alloc), nil
	case "BuckHashSys":
		return float64(m.BuckHashSys), nil
	case "Frees":
		return float64(m.Frees), nil
	case "GCCPUFraction":
		return float64(m.GCCPUFraction), nil
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
	Storage                 MemStorage
	RuntimeGaugeMetricNames []string
}

func NewMetricCollector(gaugeMetrics []string) *MetricCollector {
	tmpGaugeMetrics := make([]string, len(gaugeMetrics))
	copy(tmpGaugeMetrics, gaugeMetrics)
	return &MetricCollector{
		Storage:                 *NewMemStorage(),
		RuntimeGaugeMetricNames: tmpGaugeMetrics,
	}
}

func (m *MetricCollector) CollectMetrics() {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)
	for _, gaugeMetric := range m.RuntimeGaugeMetricNames {
		value, err := Getter(&rtm, gaugeMetric)
		if err != nil {
			panic(err)
		}
		m.Storage.UpdateGauge(gaugeMetric, value)
		m.Storage.UpdateCounter("PollCount", 1)
	}
	m.Storage.UpdateGauge("RandomValue", rand.Float64())
}

func (m *MetricCollector) DumpMetrics() *MemStorage {
	ms := NewMemStorage()
	for name, gaugeValue := range m.Storage.GaugeMetrics {
		ms.UpdateGauge(name, gaugeValue)
	}
	for name, counterValue := range m.Storage.CounterMetrics {
		ms.UpdateCounter(name, counterValue)
	}

	m.Storage.CounterMetrics["PollCount"] = 0

	return ms
}