package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/GaryShem/ya-metrics.git/internal/shared/logging"
	"github.com/GaryShem/ya-metrics.git/internal/shared/storage/metrics"
)

func (h *RepoHandler) GetGauge(w http.ResponseWriter, r *http.Request) {
	metricType := metrics.TypeGauge
	// get metric name
	metricName := chi.URLParam(r, "metricName")
	if metricName == "" {
		http.Error(w, fmt.Sprintf("%v metric type: %v", metricType, metrics.ErrInvalidMetricID), http.StatusNotFound)
		return
	}
	var valueBytes []byte
	value, err := h.repo.GetGauge(metricName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	valueBytes = []byte(fmt.Sprintf("%v", value.Value))
	if _, err = w.Write(valueBytes); err != nil {
		http.Error(w, "could not write response, contact server admins", http.StatusInternalServerError)
		return
	}
}

func (h *RepoHandler) GetCounter(w http.ResponseWriter, r *http.Request) {
	metricType := metrics.TypeCounter
	// get metric name
	metricName := chi.URLParam(r, "metricName")
	if metricName == "" {
		http.Error(w, fmt.Sprintf("%v metric type: %v", metricType, metrics.ErrInvalidMetricID), http.StatusNotFound)
		return
	}
	var valueBytes []byte
	value, err := h.repo.GetCounter(metricName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	valueBytes = []byte(fmt.Sprintf("%v", value.Value))
	if _, err = w.Write(valueBytes); err != nil {
		http.Error(w, "could not write response, contact server admins", http.StatusInternalServerError)
		return
	}
}

func (h *RepoHandler) ListMetrics(w http.ResponseWriter, _ *http.Request) {
	jsonResponse, err := json.Marshal(h.repo)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not marshal json: %v", err.Error()), http.StatusInternalServerError)
	}
	if _, err = w.Write(jsonResponse); err != nil {
		http.Error(w, "could not write response, contact server admins", http.StatusInternalServerError)
		return
	}
}

func (h *RepoHandler) GetMetricJSON(w http.ResponseWriter, r *http.Request) {
	// make sure metrics are passed via POST
	if r.Method != http.MethodPost {
		http.Error(w, "only POST is accepted", http.StatusBadRequest)
		return
	}

	// deserialize request
	metric := &metrics.Metrics{}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&metric); err != nil {
		logging.Log.Debug("error decoding request json", zap.Error(err))
		http.Error(w, "error decoding request json", http.StatusBadRequest)
		return
	}

	// get metric from repository
	err := h.repo.GetMetric(metric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// serialize updated metric structure
	response, err := json.Marshal(metric)
	if err != nil {
		logging.Log.Debug("error marshaling response", zap.Error(err))
		http.Error(w, "error marshaling response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(response); err != nil {
		http.Error(w, "could not write GetMetricJSON response, contact server admins", http.StatusInternalServerError)
		return
	}
}
