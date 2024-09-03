package grpc

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/GaryShem/ya-metrics.git/internal/shared/proto"
)

func (s *MetricsServerRepo) ListMetrics(_ context.Context, _ *emptypb.Empty) (*proto.MetricListMessage, error) {
	metrics, err := s.repo.ListMetrics()
	if err != nil {
		return nil, err
	}
	protoMetrics := make([]*proto.Metric, len(metrics))
	for i, m := range metrics {
		protoMetrics[i] = mapMetricInternalToProto(&m)
	}
	return &proto.MetricListMessage{
		Metrics: protoMetrics,
	}, nil
}
