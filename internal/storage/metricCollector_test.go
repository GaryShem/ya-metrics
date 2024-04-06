package storage

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
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
		wantErr assert.ErrorAssertionFunc
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
			wantErr: assert.NoError,
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
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Getter(tt.args.m, tt.args.metricName)
			if !tt.wantErr(t, err, fmt.Sprintf("Getter(%v, %v)", tt.args.m, tt.args.metricName)) {
				t.Errorf("Getter() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMetricCollector_CollectMetrics(t *testing.T) {
	tests := []struct {
		name                    string
		runtimeGaugeMetricNames []string
		pollCount               int64
	}{
		{
			name:                    "Test Collecting Metrics",
			runtimeGaugeMetricNames: []string{"Alloc"},
			pollCount:               1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMetricCollector(tt.runtimeGaugeMetricNames)
			m.CollectMetrics()
			for _, name := range tt.runtimeGaugeMetricNames {
				assert.NotEqual(t, 0, m.Storage.GaugeMetrics[name])
			}
			assert.Equal(t, tt.pollCount, m.Storage.CounterMetrics["PollCount"])
		})
	}
}

func TestMetricCollector_DumpMetrics(t *testing.T) {
	type fields struct {
		Storage                 MemStorage
		RuntimeGaugeMetricNames []string
	}
	type want struct {
		receivedStorage  MemStorage
		collectedStorage MemStorage
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "Test Dump Metrics",
			fields: fields{
				Storage: MemStorage{
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
				receivedStorage: MemStorage{
					GaugeMetrics: map[string]float64{
						"Alloc": 2,
					},
					CounterMetrics: map[string]int64{
						"PollCount": 1,
					},
				},
				collectedStorage: MemStorage{
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
	tests := make([]string, len(RuntimeMetrics))
	copy(tests, RuntimeMetrics)
	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			var rtm runtime.MemStats
			runtime.ReadMemStats(&rtm)
			got, err := Getter(&rtm, tt)
			assert.NoError(t, err)
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
			assert.Error(t, err)
			assert.Equal(t, 0., got)
		})
	}
}
