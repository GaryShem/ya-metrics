package grpc

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/GaryShem/ya-metrics.git/internal/shared/proto"
	"github.com/GaryShem/ya-metrics.git/internal/shared/storage/models"
)

func (s *MetricsServerRepo) UpdateGauge(_ context.Context, request *proto.UpdateGaugeRequest) (*emptypb.Empty, error) {
	if err := s.repo.UpdateGauge(request.Name, request.Value); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *MetricsServerRepo) UpdateCounter(_ context.Context, request *proto.UpdateCounterRequest) (*emptypb.Empty, error) {
	if err := s.repo.UpdateCounter(request.Name, request.Delta); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *MetricsServerRepo) UpdateMetric(_ context.Context, request *proto.MetricMessage) (*proto.MetricMessage, error) {
	metric := models.MapMetricProtoToInternal(request.Metric)
	if err := s.repo.UpdateMetric(metric); err != nil {
		return nil, err
	}
	return &proto.MetricMessage{
		Metric: models.MapMetricInternalToProto(metric),
	}, nil
}

func (s *MetricsServerRepo) UpdateBatch(_ context.Context, request *proto.MetricListMessage) (*proto.MetricListMessage, error) {
	metrics := make([]models.Metrics, 0, len(request.Metrics))
	for _, m := range request.Metrics {
		metrics = append(metrics, *models.MapMetricProtoToInternal(m))
	}
	metrics, err := s.repo.UpdateMetricBatch(metrics)
	if err != nil {
		return nil, err
	}
	result := make([]*proto.Metric, 0, len(metrics))
	for _, m := range metrics {
		result = append(result, models.MapMetricInternalToProto(&m))
	}
	return &proto.MetricListMessage{Metrics: result}, nil
}
