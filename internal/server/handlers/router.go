package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/GaryShem/ya-metrics.git/internal/server/middleware"
	"github.com/GaryShem/ya-metrics.git/internal/server/storage/postgres"
	"github.com/GaryShem/ya-metrics.git/internal/shared/logging"
	"github.com/GaryShem/ya-metrics.git/internal/shared/storage/models"
)

func MetricsRouter(ms models.Repository, dbConn *postgres.PostgreSQLStorage) (chi.Router, error) {
	if err := logging.InitializeZapLogger("Info"); err != nil {
		return nil, err
	}
	r := chi.NewRouter()
	h := NewHandler(ms)
	r.Use(middleware.RequestGzipper)
	r.Use(middleware.RequestLogger)
	r.Route(`/`, func(r chi.Router) {
		r.Get(`/ping`, dbConn.TestConnection)

		r.Route(`/update`, func(r chi.Router) {
			r.Post(`/`, h.UpdateMetricJSON)
			r.Get(`/{metricType}/{metricName}/{metricValue}`, h.UpdateMetric)
			r.Post(`/{metricType}/{metricName}/{metricValue}`, h.UpdateMetric)
			r.Post(`/{metricType}`, func(rw http.ResponseWriter, _ *http.Request) {
				http.Error(rw, "unknown metric type or no value provided", http.StatusBadRequest)
			})
		})

		r.Get(`/`, h.ListMetrics)
		r.Route(`/value`, func(r chi.Router) {
			r.Post(`/`, h.GetMetricJSON)
			r.Get(`/gauge/{metricName}`, h.GetGauge)
			r.Get(`/counter/{metricName}`, h.GetCounter)
		})
	})
	return r, nil
}
