package server

import (
	"log"
	"net/http"
	"time"

	"github.com/GaryShem/ya-metrics.git/internal/server/handlers"
	"github.com/GaryShem/ya-metrics.git/internal/server/longterm"
	"github.com/GaryShem/ya-metrics.git/internal/shared/logging"
	"github.com/GaryShem/ya-metrics.git/internal/shared/storage/memorystorage"
)

type ServerFlags struct {
	Address         *string
	StoreInterval   *int
	FileStoragePath *string
	Restore         *bool
}

func RunServer(sf *ServerFlags) error {
	if err := logging.InitializeZapLogger("Info"); err != nil {
		return err
	}
	fs := longterm.NewFileSaver(*sf.FileStoragePath, nil)
	if sf.Restore != nil && *sf.Restore {
		err := fs.LoadMetricsFile()
		if err != nil {
			return err
		}
	} else {
		fs.MS = memorystorage.NewMemStorage()
	}
	go func() {
		_ = fs.SaveMetricsFile(time.Second * time.Duration(*sf.StoreInterval))
	}()
	r, err := handlers.MetricsRouter(fs.MS)
	if err != nil {
		return err
	}
	log.Printf("Server listening on %v\n", *sf.Address)
	err = http.ListenAndServe(*sf.Address, r)
	if err != nil {
		return err
	}
	return nil
}
