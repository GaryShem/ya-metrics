package storage

type MemStorage struct {
	GaugeMetrics   map[string]float64
	CounterMetrics map[string]int64
}

type repository interface {
	UpdateGauge(metricName string, value float64)
	UpdateCounter(metricName string, value int64)
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		GaugeMetrics:   make(map[string]float64),
		CounterMetrics: make(map[string]int64),
	}
}

func (ms *MemStorage) UpdateGauge(metricName string, value float64) {
	ms.GaugeMetrics[metricName] = value
}

func (ms *MemStorage) UpdateCounter(metricName string, value int64) {
	currentValue := ms.CounterMetrics[metricName]
	ms.CounterMetrics[metricName] = currentValue + value
}

var _ repository = &MemStorage{}
