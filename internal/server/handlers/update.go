package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/GaryShem/ya-metrics.git/internal/shared/logging"
	"github.com/GaryShem/ya-metrics.git/internal/shared/storage/models"
)

func (h *RepoHandler) UpdateGauge(w http.ResponseWriter, r *http.Request) {
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
	err = h.repo.UpdateGauge(metricName, metricValue)
	if err != nil {
		http.Error(w, "could not update gauge metric", http.StatusInternalServerError)
		return
	}
	_, err = w.Write([]byte(fmt.Sprintf("metric %v updated with value %v", metricName, metricValueString)))
	if err != nil {
		http.Error(w, "could not write response, contact server admins", http.StatusInternalServerError)
		return
	}
}

func (h *RepoHandler) UpdateCounter(w http.ResponseWriter, r *http.Request) {
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
	err = h.repo.UpdateCounter(metricName, metricValue)
	if err != nil {
		http.Error(w, "could not update counter metric", http.StatusInternalServerError)
		return
	}
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
	case string(models.TypeCounter):
		h.UpdateCounter(w, r)
	case string(models.TypeGauge):
		h.UpdateGauge(w, r)
	default:
		http.Error(w, fmt.Sprintf("unknown metric type %v", metricType), http.StatusBadRequest)
	}
}

func (h *RepoHandler) UpdateMetricJSON(w http.ResponseWriter, r *http.Request) {
	// make sure metrics are passed via POST
	if r.Method != http.MethodPost {
		http.Error(w, "only POST is accepted", http.StatusBadRequest)
		return
	}

	// check content type
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "invalid content type", http.StatusNotFound)
	}

	// deserialize request
	metric := &models.Metrics{}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&metric); err != nil {
		logging.Log.Debug("error decoding request json", zap.Error(err))
		http.Error(w, "error decoding request json", http.StatusBadRequest)
	}

	// update repository
	if err := h.repo.UpdateMetric(metric); err != nil {
		logging.Log.Debug("error updating metric", zap.Error(err))
		http.Error(w, fmt.Sprintf("error updating metric, %v", err), http.StatusInternalServerError)
	}

	// serialize updated metric structure
	response, err := json.Marshal(metric)
	if err != nil {
		logging.Log.Debug("error marshaling response", zap.Error(err))
		http.Error(w, "error marshaling response", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(response); err != nil {
		http.Error(w, "could not write UpdateMetricJSON response, contact server admins", http.StatusInternalServerError)
	}
}
