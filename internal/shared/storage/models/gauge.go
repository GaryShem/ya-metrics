package models

type Gauge struct {
	Type  string  `json:"type"`
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

func NewGauge(name string, value float64) *Gauge {
	return &Gauge{
		Type:  "gauge",
		Name:  name,
		Value: value,
	}
}

func CopyGauge(original Gauge) *Gauge {
	return &Gauge{
		Type:  "gauge",
		Name:  original.Name,
		Value: original.Value,
	}
}

func (cm *Gauge) Update(value float64) {
	cm.Value = value
}

func (cm *Gauge) Reset() {
	cm.Value = 0
}
