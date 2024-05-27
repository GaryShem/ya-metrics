package models

type Repository interface {
	UpdateGauge(metricName string, value float64)
	UpdateCounter(metricName string, value int64)
	UpdateMetric(m *Metrics) error
	GetGauge(metricName string) (*Gauge, error)
	GetCounter(metricName string) (*Counter, error)
	GetMetric(m *Metrics) error
	GetGauges() map[string]*Gauge
	GetCounters() map[string]*Counter
	Ping() error
}
