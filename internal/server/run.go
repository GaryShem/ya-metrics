package server

import (
	"log"
	"net/http"

	"github.com/GaryShem/ya-metrics.git/internal/server/handlers"
	"github.com/GaryShem/ya-metrics.git/internal/shared/storage/models"
)

type ServerFlags struct {
	Address *string
}

func RunServer(sf *ServerFlags, rep models.Repository) {
	r, err := handlers.MetricsRouter(rep)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Server listening on %v\n", *sf.Address)
	err = http.ListenAndServe(*sf.Address, r)
	if err != nil {
		log.Fatal(err)
	}
}
