package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/GaryShem/ya-metrics.git/internal/server/storage/repository"
	"github.com/GaryShem/ya-metrics.git/internal/shared/logging"
)

func MetricsRouter(ms repository.Repository, middlewares ...func(http.Handler) http.Handler) (chi.Router, error) {
	if err := logging.InitializeZapLogger("Info"); err != nil {
		return nil, err
	}
	r := chi.NewRouter()
	h := NewHandler(ms)
	for _, mw := range middlewares {
		r.Use(mw)
	}

	r.Route(`/`, func(r chi.Router) {
		r.Get(`/ping`, h.Ping)
		r.Post(`/updates/`, h.UpdateMetricBatch)

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
