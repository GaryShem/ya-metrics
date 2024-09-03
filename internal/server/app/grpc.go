package app

import (
	"context"
	"errors"
	"log"
	"time"

	grpc2 "google.golang.org/grpc"

	"github.com/GaryShem/ya-metrics.git/internal/server/config"
	"github.com/GaryShem/ya-metrics.git/internal/server/handlers/grpc"
	"github.com/GaryShem/ya-metrics.git/internal/server/middleware"
	"github.com/GaryShem/ya-metrics.git/internal/server/storage/repository"
	"github.com/GaryShem/ya-metrics.git/internal/shared/proto"
)

func initGRPCServer(ctx context.Context, sf *config.ServerFlags, repo repository.Repository) error {
	if sf.GRPCAddress == "" {
		return nil
	}
	interceptors := make([]grpc2.UnaryServerInterceptor, 0)
	if sf.TrustedSubnet != "" {
		interceptor, err := middleware.NewNetworkFilterMiddleware(sf.TrustedSubnet)
		if err != nil {
			return err
		}
		interceptors = append(interceptors, interceptor.Intercept)
	}
	metrics := grpc.NewMetricsServerRepo(repo)
	server := grpc2.NewServer(grpc2.ChainUnaryInterceptor(interceptors...))
	proto.RegisterMetricsServer(server, metrics)
	go func() {
		<-ctx.Done()
		// graceful shutdown period
		shutdownCtx, shutdownStopCtx := context.WithTimeout(ctx, 10*time.Second)
		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				server.Stop()
			}
		}()

		server.GracefulStop()
		shutdownStopCtx()
	}()
	log.Printf("GRPC Server listening on %v\n", sf.GRPCAddress)

	return nil
}
