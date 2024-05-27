package main

import (
	"log"

	"github.com/GaryShem/ya-metrics.git/internal/server/app"
	"github.com/GaryShem/ya-metrics.git/internal/server/config"
)

func main() {
	sf := new(config.ServerFlags)
	config.ParseFlags(sf)
	err := app.RunServer(sf)
	if err != nil {
		log.Fatal(err)
	}
}
