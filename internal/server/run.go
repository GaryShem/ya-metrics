package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/GaryShem/ya-metrics.git/internal/server/handlers"
	"github.com/GaryShem/ya-metrics.git/internal/server/middleware"
	"github.com/GaryShem/ya-metrics.git/internal/shared/logging"
	"github.com/GaryShem/ya-metrics.git/internal/shared/storage"
)

type ServerFlags struct {
	Address *string
}

func MetricsRouter(ms storage.Repository) (chi.Router, error) {
	if err := logging.InitializeZapLogger("Info"); err != nil {
		return nil, err
	}
	r := chi.NewRouter()
	h := handlers.NewHandler(ms)
	r.Use(middleware.RequestLogger)
	r.Route(`/`, func(r chi.Router) {
		r.Route(`/update`, func(r chi.Router) {
			r.Get(`/{metricType}/{metricName}/{metricValue}`, h.UpdateMetric)
			r.Post(`/{metricType}/{metricName}/{metricValue}`, h.UpdateMetric)
			r.Post(`/`, func(rw http.ResponseWriter, _ *http.Request) {
				http.Error(rw, "no metric type provided", http.StatusBadRequest)
			})
			r.Post(`/{metricType}`, func(rw http.ResponseWriter, _ *http.Request) {
				http.Error(rw, "unknown metric type or no value provided", http.StatusBadRequest)
			})
		})

		r.Get(`/`, h.ListMetrics)
		r.Route(`/value`, func(r chi.Router) {
			r.Get(`/gauge/{metricName}`, h.GetGauge)
			r.Get(`/counter/{metricName}`, h.GetCounter)
		})
	})
	return r, nil
}

func RunServer(sf *ServerFlags, rep storage.Repository) {
	r, err := MetricsRouter(rep)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Server listening on %v\n", *sf.Address)
	err = http.ListenAndServe(*sf.Address, r)
	if err != nil {
		log.Fatal(err)
	}
}
