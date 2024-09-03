package models

import (
	"github.com/GaryShem/ya-metrics.git/internal/shared/proto"
)

func MapMetricInternalToProto(metric *Metrics) *proto.Metric {
	var g float64
	var c int64
	if metric.Value != nil {
		g = *metric.Value
	}
	if metric.Delta != nil {
		c = *metric.Delta
	}
	return &proto.Metric{
		Name:  metric.ID,
		Type:  metric.MType,
		Value: g,
		Delta: c,
	}
}

func MapMetricProtoToInternal(metric *proto.Metric) *Metrics {
	return &Metrics{
		ID:    metric.Name,
		MType: metric.Type,
		Delta: &metric.Delta,
		Value: &metric.Value,
	}
}
