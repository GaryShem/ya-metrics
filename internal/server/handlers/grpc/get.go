package proto

import (
	"context"

	"github.com/GaryShem/ya-metrics.git/internal/server/handlers/grpc/proto"
)

func (s *MetricsServer) GetGauge(_ context.Context, request *proto.NameMessage) (*proto.GetGaugeResponse, error) {
	value, err := s.repo.GetGauge(request.Name)
	if err != nil {
		return nil, err
	}
	return &proto.GetGaugeResponse{
		Value: value.Value,
	}, err
}

func (s *MetricsServer) GetCounter(_ context.Context, request *proto.NameMessage) (*proto.GetCounterResponse, error) {
	value, err := s.repo.GetCounter(request.Name)
	if err != nil {
		return nil, err
	}
	return &proto.GetCounterResponse{
		Value: value.Value,
	}, err

}

func (s *MetricsServer) GetMetric(_ context.Context, request *proto.MetricMessage) (*proto.MetricMessage, error) {
	metric := mapMetricProtoToInternal(request.Metric)
	if err := s.repo.GetMetric(metric); err != nil {
		return nil, err
	}
	return &proto.MetricMessage{Metric: mapMetricInternalToProto(metric)}, nil
}
