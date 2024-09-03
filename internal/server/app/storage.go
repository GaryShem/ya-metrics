package app

import (
	"time"

	"github.com/GaryShem/ya-metrics.git/internal/server/config"
	"github.com/GaryShem/ya-metrics.git/internal/server/longterm"
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
