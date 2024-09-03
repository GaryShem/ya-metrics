package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/GaryShem/ya-metrics.git/internal/server/config"
	"github.com/GaryShem/ya-metrics.git/internal/shared/logging"
)

func RunServer(sf *config.ServerFlags) error {
	if err := logging.InitializeZapLogger("Info"); err != nil {
		return err
	}
	repo, err := initializeStorage(sf)
	if err != nil {
		return err
	}
	logging.Log.Infoln("server starting with flags:", sf)
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	serverCtx, serverStopCtx := context.WithCancel(context.Background())
	group, ctx := errgroup.WithContext(serverCtx)
	go func() {
		sig := <-sigint
		logging.Log.Infoln("signal received:", sig)
		serverStopCtx()
	}()

	group.Go(func() error { return initHttpServer(ctx, sf, repo) })
	group.Go(func() error { return initGRPCServer(ctx, sf, repo) })

	return group.Wait()
}
