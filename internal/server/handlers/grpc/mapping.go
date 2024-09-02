package proto

import (
	"github.com/GaryShem/ya-metrics.git/internal/server/handlers/grpc/proto"
	"github.com/GaryShem/ya-metrics.git/internal/shared/storage/models"
)

func mapMetricInternalToProto(metric *models.Metrics) *proto.Metric {
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

func mapMetricProtoToInternal(metric *proto.Metric) *models.Metrics {
	return &models.Metrics{
		ID:    metric.Name,
		MType: metric.Type,
		Delta: &metric.Delta,
		Value: &metric.Value,
	}
}
