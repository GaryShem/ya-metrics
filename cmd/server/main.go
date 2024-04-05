package main

import (
	"errors"
	"github.com/GaryShem/ya-metrics.git/internal"
	"net/http"
	"strconv"
	"strings"
)

var metricStorage = internal.NewMemStorage()

func validateMetricRequest(w http.ResponseWriter, r *http.Request) error {
	// make sure metrics are passed via POST
	if r.Method != http.MethodPost {
		http.Error(w, "only POST is accepted", http.StatusBadRequest)
		return errors.New("only POST is accepted")
	}

	return nil
}

func tokenizeMetricParameters(w http.ResponseWriter, r *http.Request) ([]string, error) {
	metricParameterString := r.URL.Path
	metricParameterSlice := strings.Split(metricParameterString, "/")
	if len(metricParameterSlice) != 2 {
		http.Error(w, "invalid metric format, metric name and value expected", http.StatusNotFound)
		return nil, errors.New("invalid metric format, metric name and value expected")
	}
	return metricParameterSlice, nil
}

func updateGauge(w http.ResponseWriter, r *http.Request) {
	err := validateMetricRequest(w, r)
	if err != nil {
		return
	}

	metricParameterSlice, err := tokenizeMetricParameters(w, r)
	if err != nil {
		return
	}

	metricName := metricParameterSlice[0]
	metricValue, err := strconv.ParseFloat(metricParameterSlice[1], 64)
	if err != nil {
		http.Error(w, "invalid metric value format, expected float64", http.StatusBadRequest)
	}

	metricStorage.UpdateGauge(metricName, metricValue)
	_, err = w.Write([]byte(""))
	if err != nil {
		panic(err)
	}
}

func updateCounter(w http.ResponseWriter, r *http.Request) {
	err := validateMetricRequest(w, r)
	if err != nil {
		return
	}

	metricParameterSlice, err := tokenizeMetricParameters(w, r)
	if err != nil {
		return
	}

	metricName := metricParameterSlice[0]
	metricValue, err := strconv.ParseInt(metricParameterSlice[1], 10, 64)
	if err != nil {
		http.Error(w, "invalid metric value format, expected int64", http.StatusBadRequest)
	}

	metricStorage.UpdateCounter(metricName, metricValue)

	w.WriteHeader(200)
	_, err = w.Write([]byte(""))
	if err != nil {
		panic(err)
	}
}

func main() {
	mux := http.NewServeMux()

	mux.Handle(`/update/gauge/`, http.StripPrefix("/update/gauge/", http.HandlerFunc(updateGauge)))
	mux.Handle(`/update/counter/`, http.StripPrefix("/update/counter/", http.HandlerFunc(updateCounter)))
	mux.HandleFunc(`/`, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "", http.StatusBadRequest)
	})

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
