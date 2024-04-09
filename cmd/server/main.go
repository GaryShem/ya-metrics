package main

import (
	"github.com/GaryShem/ya-metrics.git/internal/handlers"
	"github.com/GaryShem/ya-metrics.git/internal/storage"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	ms := storage.NewMemStorage()
	mux.Handle(`/update/{metricType}/{metricName}/{metricValue}`, handlers.UpdateMetricHandler(ms))
	mux.HandleFunc(`/update/{metricType}/`, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "no metric name provided", http.StatusNotFound)
	})
	mux.HandleFunc(`/`, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "invalid method call", http.StatusBadRequest)
	})

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
