package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/GaryShem/ya-metrics.git/internal/server/handlers"
	"github.com/GaryShem/ya-metrics.git/internal/shared"
)

type ServerFlags struct {
	Address *string
}

func MetricsRouter(ms shared.Repository) chi.Router {
	r := chi.NewRouter()
	// r.Use(middleware.Logger)
	r.Route(`/`, func(r chi.Router) {
		r.Route(`/update`, func(r chi.Router) {
			r.Post(`/`, func(rw http.ResponseWriter, _ *http.Request) {
				http.Error(rw, "no metric type provided", http.StatusNotFound)
			})
			r.Route(`/{metricType}`, func(r chi.Router) {
				r.Post(`/`, func(rw http.ResponseWriter, _ *http.Request) {
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

func RunServer(sf *ServerFlags, rep shared.Repository) {
	r := MetricsRouter(rep)
	log.Printf("Server listening on %v\n", *sf.Address)
	err := http.ListenAndServe(*sf.Address, r)
	if err != nil {
		log.Fatal(err)
	}
}
