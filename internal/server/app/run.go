package app

import (
	"log"
	"net/http"
	"time"

	"github.com/GaryShem/ya-metrics.git/internal/server/config"
	"github.com/GaryShem/ya-metrics.git/internal/server/handlers"
	"github.com/GaryShem/ya-metrics.git/internal/server/longterm"
	"github.com/GaryShem/ya-metrics.git/internal/server/middleware"
	"github.com/GaryShem/ya-metrics.git/internal/server/storage/memorystorage"
	"github.com/GaryShem/ya-metrics.git/internal/server/storage/postgres"
	"github.com/GaryShem/ya-metrics.git/internal/server/storage/repository"
	"github.com/GaryShem/ya-metrics.git/internal/shared/logging"
)

func RunServer(sf *config.ServerFlags) error {
	if err := logging.InitializeZapLogger("Info"); err != nil {
		return err
	}
	logging.Log.Infoln("server starting with flags:", sf)
	var repo repository.Repository
	if sf.DBString != "" {
		logging.Log.Infoln("initializing database storage")
		dbRepo, err := postgres.NewSQLStorage(sf.DBString, true)
		if err != nil {
			return err
		}
		repo = dbRepo
	} else {
		logging.Log.Infoln("initializing memory storage")
		fs := longterm.NewFileSaver(sf.FileStoragePath, nil)
		if sf.Restore {
			err := fs.LoadMetricsFile()
			if err != nil {
				return err
			}
		} else {
			fs.MS = memorystorage.NewMemStorage()
		}
		repo = fs.MS
		go func() {
			_ = fs.SaveMetricsFile(time.Second * time.Duration(sf.StoreInterval))
		}()
	}
	middlewares := []func(http.Handler) http.Handler{}
	if sf.HashKey != "" {
		hasher := middleware.HashChecker{Key: sf.HashKey}
		middlewares = append(middlewares, hasher.Check)
	}
	middlewares = append(middlewares, middleware.RequestGzipper, middleware.RequestLogger)
	r, err := handlers.MetricsRouter(repo, middlewares...)
	if err != nil {
		return err
	}
	log.Printf("Server listening on %v\n", sf.Address)
	err = http.ListenAndServe(sf.Address, r)
	if err != nil {
		return err
	}
	return nil
}
