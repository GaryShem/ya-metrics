package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strconv"

	"github.com/go-chi/chi/v5"
)

const (
	Gauge   string = "gauge"
	Counter string = "counter"
)

func GetSupportedMetricTypes() []string {
	return []string{Gauge, Counter}
}

func (h *RepoHandler) GetMetric(w http.ResponseWriter, r *http.Request) {
	// get metric type and make sure it's an acceptable one (gauge, counter for iteration 1)
	metricType := chi.URLParam(r, "metricType")
	if !slices.Contains(GetSupportedMetricTypes(), metricType) {
		http.Error(w, fmt.Sprintf("%v metric type is not supported", metricType), http.StatusNotFound)
		return
	}
	// get metric name
	metricName := chi.URLParam(r, "metricName")
	if metricName == "" {
		http.Error(w, fmt.Sprintf("%v metric name is empty", metricType), http.StatusNotFound)
		return
	}
	var valueBytes []byte
	switch metricType {
	case Gauge:
		value, err := h.repo.GetGauge(metricName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		valueBytes = []byte(fmt.Sprintf("%v", *value))
	case Counter:
		value, err := h.repo.GetCounter(metricName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		valueBytes = []byte(fmt.Sprintf("%v", *value))
	default:
		http.Error(w, fmt.Sprintf("%v metric type is not supported", metricType), http.StatusBadRequest)
		return
	}
	if _, err := w.Write(valueBytes); err != nil {
		http.Error(w, "could not write response, contact server admins", http.StatusInternalServerError)
		return
	}
}

func (h *RepoHandler) ListMetrics(w http.ResponseWriter, r *http.Request) {
	jsonResponse, err := json.Marshal(h.repo)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not marshal json: %v", err.Error()), http.StatusInternalServerError)
	}
	if _, err = w.Write(jsonResponse); err != nil {
		http.Error(w, "could not write response, contact server admins", http.StatusInternalServerError)
		return
	}
}

func (h *RepoHandler) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	// make sure metrics are passed via POST
	if r.Method != http.MethodPost {
		http.Error(w, "only POST is accepted", http.StatusBadRequest)
		return
	}
	// get metric type and make sure it's an acceptable one (gauge, counter for iteration 1)
	metricType := chi.URLParam(r, "metricType")
	if !slices.Contains(GetSupportedMetricTypes(), metricType) {
		http.Error(w, fmt.Sprintf("%v metric type is not supported", metricType), http.StatusBadRequest)
		return
	}
	// get metric name
	metricName := chi.URLParam(r, "metricName")
	if metricName == "" {
		http.Error(w, fmt.Sprintf("%v metric name is empty", metricType), http.StatusNotFound)
		return
	}
	// get metric value and convert it into required format depending on the metric type,
	// then update corresponding metric
	metricValueString := chi.URLParam(r, "metricValue")

	switch metricType {
	case Gauge:
		metricValue, err := strconv.ParseFloat(metricValueString, 64)
		if err != nil {
			http.Error(w,
				fmt.Sprintf("%v metric value type is invalid, expected float64", metricType),
				http.StatusBadRequest)
			return
		}
		h.repo.UpdateGauge(metricName, metricValue)
	case Counter:
		metricValue, err := strconv.ParseInt(metricValueString, 10, 64)
		if err != nil {
			http.Error(w,
				fmt.Sprintf("%v metric value type is invalid, expected int64", metricType),
				http.StatusBadRequest)
			return
		}
		h.repo.UpdateCounter(metricName, metricValue)
	default:
		http.Error(w, fmt.Sprintf("%v metric type is not supported", metricType), http.StatusBadRequest)
		return
	}
	_, err := w.Write([]byte(fmt.Sprintf("metric %v updated with value %v", metricName, metricValueString)))
	if err != nil {
		http.Error(w, "could not write response, contact server admins", http.StatusInternalServerError)
		return
	}
}
