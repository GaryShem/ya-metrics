package app

import (
	"context"
	"errors"
	"log"
	"net"
	"time"

	grpc2 "google.golang.org/grpc"

	"github.com/GaryShem/ya-metrics.git/internal/server/config"
	"github.com/GaryShem/ya-metrics.git/internal/server/handlers/grpc"
	interceptors2 "github.com/GaryShem/ya-metrics.git/internal/server/interceptors"
	"github.com/GaryShem/ya-metrics.git/internal/server/storage/repository"
	"github.com/GaryShem/ya-metrics.git/internal/shared/proto"
)

func initGRPCServer(ctx context.Context, sf *config.ServerFlags, repo repository.Repository) error {
	if sf.GRPCAddress == "" {
		return nil
	}
	interceptors := make([]grpc2.UnaryServerInterceptor, 0)
	interceptors = append(interceptors, (&interceptors2.ErrorLoggingInterceptor{}).Intercept)
	if sf.TrustedSubnet != "" {
		interceptor, err := interceptors2.NewNetworkFilterInterceptor(sf.TrustedSubnet)
		if err != nil {
			return err
		}
		interceptors = append(interceptors, interceptor.Intercept)
	}
	listener, err := net.Listen("tcp", sf.GRPCAddress)
	if err != nil {
		return err
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
	if err = server.Serve(listener); err != nil {
		return err
	}
	return nil
}
