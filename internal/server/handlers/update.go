package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h *RepoHandler) UpdateGauge(w http.ResponseWriter, r *http.Request) {
	// make sure metrics are passed via POST
	if r.Method != http.MethodPost {
		http.Error(w, "only POST is accepted", http.StatusBadRequest)
		return
	}
	metricType := "gauge"
	// get metric name
	metricName := chi.URLParam(r, "metricName")
	if metricName == "" {
		http.Error(w, fmt.Sprintf("%v %v metric name is empty", metricType, metricName), http.StatusNotFound)
		return
	}
	// get metric value and convert it into required format depending on the metric type,
	// then update corresponding metric
	metricValueString := chi.URLParam(r, "metricValue")

	metricValue, err := strconv.ParseFloat(metricValueString, 64)
	if err != nil {
		http.Error(w,
			fmt.Sprintf("%v metric value type is invalid, expected float64", metricType),
			http.StatusBadRequest)
		return
	}
	h.repo.UpdateGauge(metricName, metricValue)
	_, err = w.Write([]byte(fmt.Sprintf("metric %v updated with value %v", metricName, metricValueString)))
	if err != nil {
		http.Error(w, "could not write response, contact server admins", http.StatusInternalServerError)
		return
	}
}

func (h *RepoHandler) UpdateCounter(w http.ResponseWriter, r *http.Request) {
	// make sure metrics are passed via POST
	if r.Method != http.MethodPost {
		http.Error(w, "only POST is accepted", http.StatusBadRequest)
		return
	}
	metricType := "counter"
	// get metric name
	metricName := chi.URLParam(r, "metricName")
	if metricName == "" {
		http.Error(w, fmt.Sprintf("%v %v metric name is empty", metricType, metricName), http.StatusNotFound)
		return
	}
	// get metric value and convert it into required format depending on the metric type,
	// then update corresponding metric
	metricValueString := chi.URLParam(r, "metricValue")

	metricValue, err := strconv.ParseInt(metricValueString, 10, 64)
	if err != nil {
		http.Error(w,
			fmt.Sprintf("%v metric value type is invalid, expected int64", metricType),
			http.StatusBadRequest)
		return
	}
	h.repo.UpdateCounter(metricName, metricValue)
	_, err = w.Write([]byte(fmt.Sprintf("metric %v updated with value %v", metricName, metricValueString)))
	if err != nil {
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
	metricType := chi.URLParam(r, "metricType")
	if metricType == "" {
		http.Error(w, "metric type is empty", http.StatusNotFound)
	}
	switch metricType {
	case "counter":
		h.UpdateCounter(w, r)
	case "gauge":
		h.UpdateGauge(w, r)
	default:
		http.Error(w, fmt.Sprintf("unknown metric type %v", metricType), http.StatusBadRequest)
	}
}
