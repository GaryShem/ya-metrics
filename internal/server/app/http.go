package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/GaryShem/ya-metrics.git/internal/server/config"
	http2 "github.com/GaryShem/ya-metrics.git/internal/server/handlers/http"
	"github.com/GaryShem/ya-metrics.git/internal/server/middleware"
	"github.com/GaryShem/ya-metrics.git/internal/server/storage/repository"
)

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

func initHTTPServer(ctx context.Context, sf *config.ServerFlags, repo repository.Repository) error {
	middlewares, err := initMiddlewares(sf)
	if err != nil {
		return err
	}
	//middlewares = append(middlewares, middleware.RequestLogger)
	r, err := http2.MetricsRouter(repo, false, middlewares...)
	if err != nil {
		return err
	}
	server := http.Server{
		Addr:    sf.Address,
		Handler: r,
	}
	go func() {
		<-ctx.Done()
		// graceful shutdown period
		shutdownCtx, shutdownStopCtx := context.WithTimeout(ctx, 10*time.Second)
		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				log.Fatal()
			}
		}()
		err = server.Shutdown(shutdownCtx)
		if err != nil {
			log.Print(err)
		}
		shutdownStopCtx()
	}()
	log.Printf("Server listening on %v\n", sf.Address)
	err = server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
