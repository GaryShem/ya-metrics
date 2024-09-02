package proto

import (
	"github.com/GaryShem/ya-metrics.git/internal/server/handlers/grpc/proto"
	"github.com/GaryShem/ya-metrics.git/internal/server/storage/repository"
)

type MetricsServer struct {
	proto.UnimplementedMetricsServer
	repo repository.Repository
}

var _ proto.MetricsServer = (*MetricsServer)(nil)
