package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/GaryShem/ya-metrics.git/internal/handlers"
	"github.com/GaryShem/ya-metrics.git/internal/storage"
)

func MetricsRouter(ms *storage.MemStorage) chi.Router {
	r := chi.NewRouter()
	//r.Use(middleware.Logger)
	r.Route(`/`, func(r chi.Router) {
		r.Route(`/update`, func(r chi.Router) {
			r.Post(`/`, func(rw http.ResponseWriter, r *http.Request) {
				http.Error(rw, "no metric type provided", http.StatusNotFound)
			})
			r.Route(`/{metricType}`, func(r chi.Router) {
				r.Post(`/`, func(rw http.ResponseWriter, r *http.Request) {
					http.Error(rw, "no metric name provided", http.StatusNotFound)
				})
				r.Post(`/{metricName}/{metricValue}`, handlers.UpdateMetricHandler(ms))
			})
		})
		r.Get(`/`, handlers.ListMetricsHandler(ms))
		r.Get(`/value/{metricType}/{metricName}`, handlers.FetchMetricHandler(ms))
	})
	return r
}

func main() {
	ms := storage.NewMemStorage()
	r := MetricsRouter(ms)
	sf := new(ServerFlags)
	ParseFlags(sf)
	fmt.Printf("Server listening on %v\n", *sf.address)
	err := http.ListenAndServe(*sf.address, r)
	if err != nil {
		panic(err)
	}
}
