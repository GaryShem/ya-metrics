package agent

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/GaryShem/ya-metrics.git/internal/shared"
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
		pollCount int64
	}{
		{
			name:      "Test Collecting Metrics",
			metrics:   []string{"Alloc"},
			pollCount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMetricCollector(tt.metrics)
			err := m.CollectMetrics()
			require.NoError(t, err)
			for _, name := range tt.metrics {
				assert.NotEqual(t, 0, m.Storage.GaugeMetrics[name])
			}
			assert.Equal(t, tt.pollCount, m.Storage.CounterMetrics["PollCount"])
		})
	}
}

func TestMetricCollector_DumpMetrics(t *testing.T) {
	type fields struct {
		Storage                 shared.MemStorage
		RuntimeGaugeMetricNames []string
	}
	type want struct {
		receivedStorage  shared.MemStorage
		collectedStorage shared.MemStorage
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "Test Dump Metrics",
			fields: fields{
				Storage: shared.MemStorage{
					GaugeMetrics: map[string]float64{
						"Alloc": 2,
					},
					CounterMetrics: map[string]int64{
						"PollCount": 1,
					},
				},
				RuntimeGaugeMetricNames: []string{"Alloc"},
			},
			want: want{
				receivedStorage: shared.MemStorage{
					GaugeMetrics: map[string]float64{
						"Alloc": 2,
					},
					CounterMetrics: map[string]int64{
						"PollCount": 1,
					},
				},
				collectedStorage: shared.MemStorage{
					GaugeMetrics: map[string]float64{
						"Alloc": 2,
					},
					CounterMetrics: map[string]int64{
						"PollCount": 0,
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
			assert.Equalf(t, tt.want.receivedStorage, *m.DumpMetrics(), "received storage")
			assert.Equalf(t, tt.want.collectedStorage, m.Storage, "updated collector storage")
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
				assert.Equal(t, 0., got)
			}
		})
	}
}
