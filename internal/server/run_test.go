package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/GaryShem/ya-metrics.git/internal/shared"
)

func TestMetricHandler(t *testing.T) {
	ts := httptest.NewServer(MetricsRouter(
		&shared.MemStorage{
			GaugeMetrics:   map[string]float64{},
			CounterMetrics: map[string]int64{},
		},
	))
	type testRequest struct {
		method string
		url    string
		code   int
		want   string
	}
	tests := []struct {
		name          string
		updateRequest []testRequest
		getRequest    []testRequest
	}{
		{
			name: "Gauge foo",
			updateRequest: []testRequest{
				{
					method: http.MethodPost,
					url:    "/update/gauge/foo/2",
					code:   200,
					want:   "",
				},
			},
			getRequest: []testRequest{
				{
					method: http.MethodGet,
					url:    "/value/gauge/foo",
					code:   200,
					want:   "2",
				},
			},
		},
		{
			name: "Gauge bar/barr invalid read",
			updateRequest: []testRequest{
				{
					method: http.MethodPost,
					url:    "/update/gauge/bar/2",
					code:   200,
					want:   "",
				},
			},
			getRequest: []testRequest{
				{
					method: http.MethodGet,
					url:    "/value/gauge/barr",
					code:   404,
					want:   "2",
				},
			},
		},

		{
			name: "Counter foo",
			updateRequest: []testRequest{
				{
					method: http.MethodPost,
					url:    "/update/counter/foo/2",
					code:   200,
					want:   "",
				},
			},
			getRequest: []testRequest{
				{
					method: http.MethodGet,
					url:    "/value/counter/foo",
					code:   200,
					want:   "2",
				},
			},
		},
		{
			name: "Counter bar/barr invalid read",
			updateRequest: []testRequest{
				{
					method: http.MethodPost,
					url:    "/update/counter/bar/2",
					code:   200,
					want:   "",
				},
				{
					method: http.MethodPost,
					url:    "/update/counter/bar/2",
					code:   200,
					want:   "",
				},
			},
			getRequest: []testRequest{
				{
					method: http.MethodGet,
					url:    "/value/counter/barr",
					code:   404,
					want:   "",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := resty.New()
			for _, req := range tt.updateRequest {
				url := ts.URL + req.url
				var response *resty.Response
				var err error
				switch req.method {
				case http.MethodPost:
					response, err = client.R().Post(url)
				case http.MethodGet:
					response, err = client.R().Get(url)
				default:
					panic("invalid http method")
				}
				require.NoError(t, err)
				assert.Equal(t, req.code, response.StatusCode())
			}
			for _, req := range tt.getRequest {
				url := ts.URL + req.url
				var response *resty.Response
				var err error
				switch req.method {
				case http.MethodPost:
					response, err = client.R().Post(url)
				case http.MethodGet:
					response, err = client.R().Get(url)
				default:
					panic("invalid http method")
				}
				require.NoError(t, err)
				assert.Equal(t, req.code, response.StatusCode())
				if req.code == 200 {
					assert.Equal(t, req.want, string(response.Body()))
				}
			}
		})
	}
}
