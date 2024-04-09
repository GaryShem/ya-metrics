package handlers

import (
	"github.com/GaryShem/ya-metrics.git/internal/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestUpdateMetricHandler(t *testing.T) {
	type want struct {
		status int
		data   storage.MemStorage
	}
	type updateRequest struct {
		target     string
		metricType string
		name       string
		value      string
	}
	tests := []struct {
		name    string
		data    *storage.MemStorage
		request updateRequest
		want    want
	}{
		{
			name: "Create Gauge Metric",
			data: storage.NewMemStorage(),
			request: updateRequest{
				target:     "/update/gauge/foo/200",
				metricType: "gauge",
				name:       "foo",
				value:      "200",
			},
			want: want{
				status: http.StatusOK,
				data: storage.MemStorage{
					GaugeMetrics: map[string]float64{
						"foo": 200,
					},
					CounterMetrics: map[string]int64{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := storage.NewMemStorage()
			umh := UpdateMetricHandler(ms)
			req := httptest.NewRequest(http.MethodPost, tt.request.target, nil)
			req.SetPathValue(`metricType`, tt.request.metricType)
			req.SetPathValue(`metricName`, tt.request.name)
			req.SetPathValue(`metricValue`, tt.request.value)
			w := httptest.NewRecorder()
			umh.ServeHTTP(w, req)
			t.Log(w.Body.String())
			assert.Equal(t, tt.want.status, w.Code)
			assert.Equal(t, true, reflect.DeepEqual(tt.want.data, *ms))
		})
	}
}
