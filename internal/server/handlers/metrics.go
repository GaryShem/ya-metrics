package handlers

const (
	Gauge   string = "gauge"
	Counter string = "counter"
)

func GetSupportedMetricTypes() []string {
	return []string{Gauge, Counter}
}
