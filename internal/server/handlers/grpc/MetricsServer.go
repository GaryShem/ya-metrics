package grpc

import (
	"github.com/GaryShem/ya-metrics.git/internal/server/storage/repository"
	"github.com/GaryShem/ya-metrics.git/internal/shared/proto"
)

type MetricsServerRepo struct {
	proto.UnimplementedMetricsServer
	repo repository.Repository
}

func NewMetricsServerRepo(repo repository.Repository) *MetricsServerRepo {
	return &MetricsServerRepo{repo: repo}
}

var _ proto.MetricsServer = (*MetricsServerRepo)(nil)
