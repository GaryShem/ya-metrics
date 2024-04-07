package shared

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestMemStorage_GetGauge(t *testing.T) {
	type data struct {
		gaugeMetrics   map[string]float64
		counterMetrics map[string]int64
	}
	tests := []struct {
		name    string
		data    data
		metric  string
		want    float64
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "Get Valid Gauge Metric",
			data: data{
				gaugeMetrics: map[string]float64{
					"foo": 3.14,
				},
				counterMetrics: map[string]int64{},
			},
			metric:  "foo",
			want:    3.14,
			wantErr: require.NoError,
		},
		{
			name: "Get Invalid Gauge Metric",
			data: data{
				gaugeMetrics: map[string]float64{
					"foo": 3.14,
				},
				counterMetrics: map[string]int64{},
			},
			metric:  "bar",
			want:    0,
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
		gaugeMetrics   map[string]float64
		counterMetrics map[string]int64
	}
	tests := []struct {
		name    string
		data    data
		metric  string
		want    int64
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Get Valid Gauge Metric",
			data: data{
				gaugeMetrics: map[string]float64{},
				counterMetrics: map[string]int64{
					"foo": 3,
				},
			},
			metric:  "foo",
			want:    3,
			wantErr: assert.NoError,
		},
		{
			name: "Get Invalid Gauge Metric",
			data: data{
				gaugeMetrics: map[string]float64{},
				counterMetrics: map[string]int64{
					"foo": 3,
				},
			},
			metric:  "bar",
			want:    0,
			wantErr: assert.Error,
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
				GaugeMetrics:   map[string]float64{},
				CounterMetrics: map[string]int64{},
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
		GaugeMetrics   map[string]float64
		CounterMetrics map[string]int64
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]int64
	}{
		{
			name: "Get Counters Test",
			fields: fields{
				GaugeMetrics: map[string]float64{},
				CounterMetrics: map[string]int64{
					"foo": 42,
				},
			},
			want: map[string]int64{"foo": 42},
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
		GaugeMetrics   map[string]float64
		CounterMetrics map[string]int64
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]float64
	}{
		{
			name: "Get Counters Test",
			fields: fields{
				GaugeMetrics: map[string]float64{
					"foo": 42,
				},
				CounterMetrics: map[string]int64{},
			},
			want: map[string]float64{"foo": 42},
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
