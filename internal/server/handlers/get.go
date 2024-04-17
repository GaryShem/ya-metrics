package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *RepoHandler) GetGauge(w http.ResponseWriter, r *http.Request) {
	metricType := "gauge"
	// get metric name
	metricName := chi.URLParam(r, "metricName")
	if metricName == "" {
		http.Error(w, fmt.Sprintf("%v %v metric name is empty", metricType, metricName), http.StatusNotFound)
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
	metricType := "counter"
	// get metric name
	metricName := chi.URLParam(r, "metricName")
	if metricName == "" {
		http.Error(w, fmt.Sprintf("%v %v metric name is empty", metricType, metricName), http.StatusNotFound)
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
