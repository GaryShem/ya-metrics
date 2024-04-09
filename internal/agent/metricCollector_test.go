package agent

import (
	"encoding/json"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/GaryShem/ya-metrics.git/internal/shared/storage"
	"github.com/GaryShem/ya-metrics.git/internal/shared/storage/metrics"
)

func TestGetter(t *testing.T) {
	type args struct {
		m          *runtime.MemStats
		metricName string
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "Allowed Metric",
			args: args{
				m: &runtime.MemStats{
					Alloc: 1,
				},
				metricName: "Alloc",
			},
			want:    1,
			wantErr: require.NoError,
		},
		{
			name: "Not Allowed Metric",
			args: args{
				m: &runtime.MemStats{
					Alloc: 1,
				},
				metricName: "wololo",
			},
			want:    0,
			wantErr: require.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Getter(tt.args.m, tt.args.metricName)
			tt.wantErr(t, err)
			if err == nil {
				assert.InEpsilon(t, tt.want, got, 0.001)
			}
		})
	}
}

func TestMetricCollector_CollectMetrics(t *testing.T) {
	tests := []struct {
		name      string
		metrics   []string
		pollCount *metrics.Counter
	}{
		{
			name:      "Test Collecting Metrics",
			metrics:   []string{"Alloc"},
			pollCount: metrics.NewCounter("PollCount", 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMetricCollector(tt.metrics)
			err := m.CollectMetrics()
			require.NoError(t, err)
			for _, name := range tt.metrics {
				assert.NotEqual(t, 0, m.Storage.GetGauges()[name])
			}
			pollCount, err := m.Storage.GetCounter("PollCount")
			require.NoError(t, err)
			assert.Equal(t, tt.pollCount, pollCount)
		})
	}
}

func TestMetricCollector_DumpMetrics(t *testing.T) {
	type fields struct {
		Storage                 *storage.MemStorage
		RuntimeGaugeMetricNames []string
	}
	type want struct {
		receivedStorage  *storage.MemStorage
		collectedStorage *storage.MemStorage
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "Test Dump Metrics",
			fields: fields{
				Storage: &storage.MemStorage{
					GaugeMetrics: map[string]*metrics.Gauge{
						"Alloc": metrics.NewGauge("Alloc", 2),
					},
					CounterMetrics: map[string]*metrics.Counter{
						"PollCount": metrics.NewCounter("PollCount", 1),
					},
				},
				RuntimeGaugeMetricNames: []string{"Alloc"},
			},
			want: want{
				receivedStorage: &storage.MemStorage{
					GaugeMetrics: map[string]*metrics.Gauge{
						"Alloc": metrics.NewGauge("Alloc", 2),
					},
					CounterMetrics: map[string]*metrics.Counter{
						"PollCount": metrics.NewCounter("PollCount", 1),
					},
				},
				collectedStorage: &storage.MemStorage{
					GaugeMetrics: map[string]*metrics.Gauge{
						"Alloc": metrics.NewGauge("Alloc", 2),
					},
					CounterMetrics: map[string]*metrics.Counter{
						"PollCount": metrics.NewCounter("PollCount", 0),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MetricCollector{
				Storage:                 tt.fields.Storage,
				RuntimeGaugeMetricNames: tt.fields.RuntimeGaugeMetricNames,
			}
			metricDump, err := m.DumpMetrics()
			require.NoError(t, err)
			wantRcvd, err := json.Marshal(tt.want.receivedStorage)
			require.NoError(t, err)
			gotRcvd, err := json.Marshal(metricDump)
			require.NoError(t, err)
			assert.Equal(t, string(wantRcvd), string(gotRcvd))

			wantCollected, err := json.Marshal(tt.want.collectedStorage)
			require.NoError(t, err)
			gotCollected, err := json.Marshal(m.Storage)
			require.NoError(t, err)
			assert.Equal(t, string(wantCollected), string(gotCollected))
		})
	}
}

func TestGetter_ValidMetrics(t *testing.T) {
	runtimeMetrics := SupportedRuntimeMetrics()
	tests := make([]string, len(runtimeMetrics))
	copy(tests, runtimeMetrics)
	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			var rtm runtime.MemStats
			runtime.ReadMemStats(&rtm)
			got, err := Getter(&rtm, tt)
			require.NoError(t, err)
			assert.NotEqual(t, 0, got)
		})
	}
}

func TestGetter_InvalidMetrics(t *testing.T) {
	tests := []string{
		"wololo",
	}
	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			var rtm runtime.MemStats
			runtime.ReadMemStats(&rtm)
			got, err := Getter(&rtm, tt)
			require.Error(t, err)
			if err == nil {
				assert.Equal(t, 0, got)
			}
		})
	}
}
