package grpc

import (
	"context"

	"github.com/GaryShem/ya-metrics.git/internal/shared/proto"
	"github.com/GaryShem/ya-metrics.git/internal/shared/storage/models"
)

func (s *MetricsServerRepo) GetGauge(_ context.Context, request *proto.NameMessage) (*proto.GetGaugeResponse, error) {
	value, err := s.repo.GetGauge(request.Name)
	if err != nil {
		return nil, err
	}
	return &proto.GetGaugeResponse{
		Value: value.Value,
	}, err
}

func (s *MetricsServerRepo) GetCounter(_ context.Context, request *proto.NameMessage) (*proto.GetCounterResponse, error) {
	value, err := s.repo.GetCounter(request.Name)
	if err != nil {
		return nil, err
	}
	return &proto.GetCounterResponse{
		Value: value.Value,
	}, err

}

func (s *MetricsServerRepo) GetMetric(_ context.Context, request *proto.MetricMessage) (*proto.MetricMessage, error) {
	metric := models.MapMetricProtoToInternal(request.Metric)
	if err := s.repo.GetMetric(metric); err != nil {
		return nil, err
	}
	return &proto.MetricMessage{Metric: models.MapMetricInternalToProto(metric)}, nil
}