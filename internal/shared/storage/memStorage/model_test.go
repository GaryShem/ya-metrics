package memstorage

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/GaryShem/ya-metrics.git/internal/shared/storage/metrics"
)

func TestNewMemStorage(t *testing.T) {
	tests := []struct {
		name string
		want *MemStorage
	}{
		{
			name: "Test New Mem Storage",
			want: &MemStorage{
				GaugeMetrics:   map[string]*metrics.Gauge{},
				CounterMetrics: map[string]*metrics.Counter{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewMemStorage(), "NewMemStorage()")
		})
	}
}

func TestMemStorage_GetCounters(t *testing.T) {
	type fields struct {
		GaugeMetrics   map[string]*metrics.Gauge
		CounterMetrics map[string]*metrics.Counter
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]*metrics.Counter
	}{
		{
			name: "Get Counters Test",
			fields: fields{
				GaugeMetrics: map[string]*metrics.Gauge{},
				CounterMetrics: map[string]*metrics.Counter{
					"foo": metrics.NewCounter("foo", 42),
				},
			},
			want: map[string]*metrics.Counter{"foo": metrics.NewCounter("foo", 42)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &MemStorage{
				GaugeMetrics:   tt.fields.GaugeMetrics,
				CounterMetrics: tt.fields.CounterMetrics,
			}
			assert.Equalf(t, tt.want, ms.GetCounters(), "GetCounters()")
		})
	}
}

func TestMemStorage_GetGauges(t *testing.T) {
	type fields struct {
		GaugeMetrics   map[string]*metrics.Gauge
		CounterMetrics map[string]*metrics.Counter
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]*metrics.Gauge
	}{
		{
			name: "Get Counters Test",
			fields: fields{
				GaugeMetrics: map[string]*metrics.Gauge{
					"foo": metrics.NewGauge("foo", 42),
				},
				CounterMetrics: map[string]*metrics.Counter{},
			},
			want: map[string]*metrics.Gauge{"foo": metrics.NewGauge("foo", 42)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &MemStorage{
				GaugeMetrics:   tt.fields.GaugeMetrics,
				CounterMetrics: tt.fields.CounterMetrics,
			}
			assert.Equalf(t, tt.want, ms.GetGauges(), "GetGauges()")
		})
	}
}
