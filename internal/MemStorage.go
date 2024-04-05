package internal

type MemStorage struct {
	gaugeMetrics   map[string]float64
	counterMetrics map[string]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gaugeMetrics:   make(map[string]float64),
		counterMetrics: make(map[string]int64),
	}
}

func (ms *MemStorage) UpdateGauge(metricName string, value float64) {
	ms.gaugeMetrics[metricName] = value
}

func (ms *MemStorage) UpdateCounter(metricName string, value int64) {
	currentValue, _ := ms.counterMetrics[metricName]
	ms.counterMetrics[metricName] = currentValue + value
}
