package shared

type Repository interface {
	UpdateGauge(metricName string, value float64)
	UpdateCounter(metricName string, value int64)
	GetGauge(metricName string) (float64, error)
	GetCounter(metricName string) (int64, error)
	GetGauges() map[string]float64
	GetCounters() map[string]int64
}
