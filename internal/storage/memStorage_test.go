package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemStorage_UpdateGauge(t *testing.T) {
	type fields struct {
		gaugeMetrics   map[string]float64
		counterMetrics map[string]int64
	}
	type args struct {
		metricName string
		value      float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   fields
	}{
		{
			name: "Create gauge metric 1",
			fields: fields{
				gaugeMetrics:   map[string]float64{},
				counterMetrics: map[string]int64{},
			},
			args: args{
				metricName: "foo",
				value:      100500,
			},
			want: fields{
				gaugeMetrics: map[string]float64{
					"foo": 100500,
				},
				counterMetrics: map[string]int64{},
			},
		},
		{
			name: "Update gauge metric 1",
			fields: fields{
				gaugeMetrics: map[string]float64{
					"foo": 100500,
					"bar": 3,
				},
				counterMetrics: map[string]int64{},
			},
			args: args{
				metricName: "foo",
				value:      100501,
			},
			want: fields{
				gaugeMetrics: map[string]float64{
					"foo": 100501,
					"bar": 3,
				},
				counterMetrics: map[string]int64{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &MemStorage{
				GaugeMetrics:   tt.fields.gaugeMetrics,
				CounterMetrics: tt.fields.counterMetrics,
			}
			ms.UpdateGauge(tt.args.metricName, tt.args.value)
			assert.Equal(t, tt.want.gaugeMetrics, ms.GaugeMetrics)
			assert.Equal(t, tt.want.counterMetrics, ms.CounterMetrics)
		})
	}
}

func TestMemStorage_UpdateCounter(t *testing.T) {
	type fields struct {
		gaugeMetrics   map[string]float64
		counterMetrics map[string]int64
	}
	type args struct {
		metricName string
		value      int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   fields
	}{
		{
			name: "Create gauge metric 1",
			fields: fields{
				gaugeMetrics:   map[string]float64{},
				counterMetrics: map[string]int64{},
			},
			args: args{
				metricName: "foo",
				value:      5,
			},
			want: fields{
				gaugeMetrics: map[string]float64{},
				counterMetrics: map[string]int64{
					"foo": 5,
				},
			},
		},
		{
			name: "Update gauge metric 1",
			fields: fields{
				gaugeMetrics: map[string]float64{},
				counterMetrics: map[string]int64{
					"foo": 5,
					"bar": 3,
				},
			},
			args: args{
				metricName: "foo",
				value:      5,
			},
			want: fields{
				gaugeMetrics: map[string]float64{},
				counterMetrics: map[string]int64{
					"foo": 10,
					"bar": 3,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := MemStorage{
				GaugeMetrics:   tt.fields.gaugeMetrics,
				CounterMetrics: tt.fields.counterMetrics,
			}
			ms.UpdateCounter(tt.args.metricName, tt.args.value)
			assert.Equal(t, tt.want.gaugeMetrics, ms.GaugeMetrics)
			assert.Equal(t, tt.want.counterMetrics, ms.CounterMetrics)
		})
	}
}
