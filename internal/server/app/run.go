package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/GaryShem/ya-metrics.git/internal/server/config"
	http2 "github.com/GaryShem/ya-metrics.git/internal/server/handlers/http"
	"github.com/GaryShem/ya-metrics.git/internal/server/longterm"
	"github.com/GaryShem/ya-metrics.git/internal/server/middleware"
	"github.com/GaryShem/ya-metrics.git/internal/server/storage/memorystorage"
	"github.com/GaryShem/ya-metrics.git/internal/server/storage/postgres"
	"github.com/GaryShem/ya-metrics.git/internal/server/storage/repository"
	"github.com/GaryShem/ya-metrics.git/internal/shared/logging"
)

func initializeStorage(sf *config.ServerFlags) (repository.Repository, error) {
	var repo repository.Repository
	if sf.DBString != "" {
		logging.Log.Infoln("initializing database storage")
		dbRepo, err := postgres.NewSQLStorage(sf.DBString, true)
		if err != nil {
			return nil, err
		}
		repo = dbRepo
	} else {
		logging.Log.Infoln("initializing memory storage")
		fs := longterm.NewFileSaver(sf.FileStoragePath, nil)
		if sf.Restore {
			err := fs.LoadMetricsFile()
			if err != nil {
				return nil, err
			}
		} else {
			fs.MS = memorystorage.NewMemStorage()
		}
		repo = fs.MS
		go func() {
			_ = fs.SaveMetricsFile(time.Second * time.Duration(sf.StoreInterval))
		}()
	}
	return repo, nil
}

func initMiddlewares(sf *config.ServerFlags) ([]func(http.Handler) http.Handler, error) {
	middlewares := []func(http.Handler) http.Handler{}
	if sf.TrustedSubnet != "" {
		networkFilter, err := middleware.NewNetworkFilterMiddleware(sf.TrustedSubnet)
		if err != nil {
			return nil, err
		}
		middlewares = append(middlewares, networkFilter.Validate)
	}
	if sf.CryptoKey != "" {
		decryptor, err := middleware.NewEncryptionMiddleware(sf.CryptoKey)
		if err != nil {
			return nil, err
		}
		middlewares = append(middlewares, decryptor.Decrypt)
	}
	if sf.HashKey != "" {
		hasher := middleware.HashChecker{Key: sf.HashKey}
		middlewares = append(middlewares, hasher.Check)
	}
	middlewares = append(middlewares, middleware.RequestGzipper)
	return middlewares, nil
}

func RunServer(sf *config.ServerFlags) error {
	if err := logging.InitializeZapLogger("Info"); err != nil {
		return err
	}
	logging.Log.Infoln("server starting with flags:", sf)
	repo, err := initializeStorage(sf)
	if err != nil {
		return err
	}
	middlewares, err := initMiddlewares(sf)
	if err != nil {
		return err
	}
	//middlewares = append(middlewares, middleware.RequestLogger)
	r, err := http2.MetricsRouter(repo, false, middlewares...)
	if err != nil {
		return err
	}
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	server := http.Server{
		Addr:    sf.Address,
		Handler: r,
	}
	serverCtx, serverStopCtx := context.WithCancel(context.Background())
	go func() {
		<-sigint
		// graceful shutdown period
		shutdownCtx, shutdownStopCtx := context.WithTimeout(serverCtx, 10*time.Second)
		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				log.Fatal()
			}
		}()
		err = server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		shutdownStopCtx()
		serverStopCtx()
	}()
	log.Printf("Server listening on %v\n", sf.Address)
	err = server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	<-serverCtx.Done()
	return nil
}
