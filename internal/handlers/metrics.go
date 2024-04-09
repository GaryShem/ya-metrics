package handlers

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"

	"github.com/GaryShem/ya-metrics.git/internal/storage"
)

const (
	Gauge   string = "gauge"
	Counter string = "counter"
)

var supportedMetricTypes = []string{Gauge, Counter}

func UpdateMetricHandler(ms *storage.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// make sure metrics are passed via POST
		if r.Method != http.MethodPost {
			http.Error(w, "only POST is accepted", http.StatusBadRequest)
			return
		}
		// get metric type and make sure it's an acceptable one (gauge, counter for iteration 1)
		metricType := r.PathValue("metricType")
		if !slices.Contains(supportedMetricTypes, metricType) {
			http.Error(w, fmt.Sprintf("%v metric type is not supported", metricType), http.StatusBadRequest)
			return
		}
		// get metric name
		metricName := r.PathValue("metricName")
		if metricName == "" {
			http.Error(w, fmt.Sprintf("%v metric name is empty", metricType), http.StatusNotFound)
		}
		// get metric value and convert it into required format depending on the metric type,
		// then update corresponding metric
		metricValueString := r.PathValue("metricValue")
		if metricType == Gauge {
			metricValue, err := strconv.ParseFloat(metricValueString, 64)
			if err != nil {
				http.Error(w,
					fmt.Sprintf("%v metric value type is invalid, expected float64", metricType),
					http.StatusBadRequest)
				return
			}
			ms.UpdateGauge(metricName, metricValue)
		} else if metricType == Counter {
			metricValue, err := strconv.ParseInt(metricValueString, 10, 64)
			if err != nil {
				http.Error(w,
					fmt.Sprintf("%v metric value type is invalid, expected int64", metricType),
					http.StatusBadRequest)
				return
			}
			ms.UpdateCounter(metricName, metricValue)
		}

		_, err := w.Write([]byte(""))
		if err != nil {
			panic(err)
		}
	}
}
