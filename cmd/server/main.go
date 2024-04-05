package main

import (
	"fmt"
	"github.com/GaryShem/ya-metrics.git/internal"
	"net/http"
	"slices"
	"strconv"
)

var metricStorage = internal.NewMemStorage()

const (
	Gauge   string = "gauge"
	Counter string = "counter"
)

var supportedMetricTypes = []string{Gauge, Counter}

func updateMetric(w http.ResponseWriter, r *http.Request) {
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
		metricStorage.UpdateGauge(metricName, metricValue)
	} else if metricType == Counter {
		metricValue, err := strconv.ParseInt(metricValueString, 10, 64)
		if err != nil {
			http.Error(w,
				fmt.Sprintf("%v metric value type is invalid, expected int64", metricType),
				http.StatusBadRequest)
			return
		}
		metricStorage.UpdateCounter(metricName, metricValue)
	}

	w.WriteHeader(200)
	_, err := w.Write([]byte(""))
	if err != nil {
		panic(err)
	}
}

func main() {
	mux := http.NewServeMux()

	mux.Handle(`/update/{metricType}/{metricName}/{metricValue}`, http.StripPrefix("/update/", http.HandlerFunc(updateMetric)))
	mux.HandleFunc(`/`, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "", http.StatusBadRequest)
	})

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
