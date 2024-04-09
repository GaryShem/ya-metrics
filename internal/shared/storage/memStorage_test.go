package storage

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/GaryShem/ya-metrics.git/internal/shared/storage/metrics"
)

func TestMemStorage_UpdateGauge(t *testing.T) {
	type fields struct {
		gaugeMetrics   map[string]*metrics.Gauge
		counterMetrics map[string]*metrics.Counter
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
				gaugeMetrics:   map[string]*metrics.Gauge{},
				counterMetrics: map[string]*metrics.Counter{},
			},
			args: args{
				metricName: "foo",
				value:      100500,
			},
			want: fields{
				gaugeMetrics: map[string]*metrics.Gauge{
					"foo": metrics.NewGauge("foo", 100500),
				},
				counterMetrics: map[string]*metrics.Counter{},
			},
		},
		{
			name: "Update gauge metric 1",
			fields: fields{
				gaugeMetrics: map[string]*metrics.Gauge{
					"foo": metrics.NewGauge("foo", 3),
				},
				counterMetrics: map[string]*metrics.Counter{},
			},
			args: args{
				metricName: "foo",
				value:      100501,
			},
			want: fields{
				gaugeMetrics: map[string]*metrics.Gauge{
					"foo": metrics.NewGauge("foo", 100501),
					//"bar": metrics.NewGauge("bar", 3),
				},
				counterMetrics: map[string]*metrics.Counter{},
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
			wantGauge, err := json.Marshal(tt.want.gaugeMetrics)
			require.NoError(t, err)
			gotGauge, err := json.Marshal(ms.GetGauges())
			require.NoError(t, err)
			assert.Equal(t, string(wantGauge), string(gotGauge))

			wantCounter, err := json.Marshal(tt.want.counterMetrics)
			require.NoError(t, err)
			gotCounter, err := json.Marshal(ms.GetCounters())
			require.NoError(t, err)
			assert.Equal(t, string(wantCounter), string(gotCounter))
		})
	}
}

func TestMemStorage_UpdateCounter(t *testing.T) {
	type fields struct {
		gaugeMetrics   map[string]*metrics.Gauge
		counterMetrics map[string]*metrics.Counter
	}
	type args struct {
		metricName string
		value      int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   MemStorage
	}{
		{
			name: "Create gauge metric 1",
			fields: fields{
				gaugeMetrics:   map[string]*metrics.Gauge{},
				counterMetrics: map[string]*metrics.Counter{},
			},
			args: args{
				metricName: "foo",
				value:      5,
			},
			want: MemStorage{
				GaugeMetrics: map[string]*metrics.Gauge{},
				CounterMetrics: map[string]*metrics.Counter{
					"foo": metrics.NewCounter("foo", 5),
				},
			},
		},
		{
			name: "Update gauge metric 1",
			fields: fields{
				gaugeMetrics: map[string]*metrics.Gauge{},
				counterMetrics: map[string]*metrics.Counter{
					"foo": metrics.NewCounter("foo", 5),
					"bar": metrics.NewCounter("bar", 3),
				},
			},
			args: args{
				metricName: "foo",
				value:      5,
			},
			want: MemStorage{
				GaugeMetrics: map[string]*metrics.Gauge{},
				CounterMetrics: map[string]*metrics.Counter{
					"foo": metrics.NewCounter("foo", 10),
					"bar": metrics.NewCounter("bar", 3),
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
			wantJSON, err := json.Marshal(tt.want)
			require.NoError(t, err)
			gotJSON, err := json.Marshal(ms)
			require.NoError(t, err)
			assert.Equal(t, string(wantJSON), string(gotJSON))
		})
	}
}

func TestMemStorage_GetGauge(t *testing.T) {
	type data struct {
		gaugeMetrics   map[string]*metrics.Gauge
		counterMetrics map[string]*metrics.Counter
	}
	tests := []struct {
		name    string
		data    data
		metric  string
		want    *metrics.Gauge
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "Get Valid Gauge Metric",
			data: data{
				gaugeMetrics: map[string]*metrics.Gauge{
					"foo": metrics.NewGauge("foo", 3.14),
				},
				counterMetrics: map[string]*metrics.Counter{},
			},
			metric:  "foo",
			want:    metrics.NewGauge("foo", 3.14),
			wantErr: require.NoError,
		},
		{
			name: "Get Invalid Gauge Metric",
			data: data{
				gaugeMetrics: map[string]*metrics.Gauge{
					"foo": metrics.NewGauge("foo", 3.14),
				},
				counterMetrics: map[string]*metrics.Counter{},
			},
			metric:  "bar",
			want:    nil,
			wantErr: require.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &MemStorage{
				GaugeMetrics:   tt.data.gaugeMetrics,
				CounterMetrics: tt.data.counterMetrics,
			}
			got, err := ms.GetGauge(tt.metric)
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMemStorage_GetCounter(t *testing.T) {
	type data struct {
		gaugeMetrics   map[string]*metrics.Gauge
		counterMetrics map[string]*metrics.Counter
	}
	tests := []struct {
		name    string
		data    data
		metric  string
		want    *metrics.Counter
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "Get Valid Gauge Metric",
			data: data{
				gaugeMetrics: map[string]*metrics.Gauge{},
				counterMetrics: map[string]*metrics.Counter{
					"foo": metrics.NewCounter("foo", 3),
				},
			},
			metric:  "foo",
			want:    metrics.NewCounter("foo", 3),
			wantErr: require.NoError,
		},
		{
			name: "Get Invalid Gauge Metric",
			data: data{
				gaugeMetrics: map[string]*metrics.Gauge{},
				counterMetrics: map[string]*metrics.Counter{
					"foo": metrics.NewCounter("foo", 3),
				},
			},
			metric:  "bar",
			want:    nil,
			wantErr: require.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &MemStorage{
				GaugeMetrics:   tt.data.gaugeMetrics,
				CounterMetrics: tt.data.counterMetrics,
			}
			got, err := ms.GetCounter(tt.metric)
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

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
