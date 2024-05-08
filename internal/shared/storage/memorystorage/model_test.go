package memorystorage

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/GaryShem/ya-metrics.git/internal/shared/storage/models"
)

func TestNewMemStorage(t *testing.T) {
	tests := []struct {
		name string
		want *MemStorage
	}{
		{
			name: "Test New Mem Storage",
			want: &MemStorage{
				GaugeMetrics:   map[string]*models.Gauge{},
				CounterMetrics: map[string]*models.Counter{},
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
		GaugeMetrics   map[string]*models.Gauge
		CounterMetrics map[string]*models.Counter
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]*models.Counter
	}{
		{
			name: "Get Counters Test",
			fields: fields{
				GaugeMetrics: map[string]*models.Gauge{},
				CounterMetrics: map[string]*models.Counter{
					"foo": models.NewCounter("foo", 42),
				},
			},
			want: map[string]*models.Counter{"foo": models.NewCounter("foo", 42)},
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
		GaugeMetrics   map[string]*models.Gauge
		CounterMetrics map[string]*models.Counter
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]*models.Gauge
	}{
		{
			name: "Get Counters Test",
			fields: fields{
				GaugeMetrics: map[string]*models.Gauge{
					"foo": models.NewGauge("foo", 42),
				},
				CounterMetrics: map[string]*models.Counter{},
			},
			want: map[string]*models.Gauge{"foo": models.NewGauge("foo", 42)},
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
