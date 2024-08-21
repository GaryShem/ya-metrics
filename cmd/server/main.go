package main

import (
	"log"

	"github.com/GaryShem/ya-metrics.git/internal/server/app"
	"github.com/GaryShem/ya-metrics.git/internal/server/config"
	"github.com/GaryShem/ya-metrics.git/internal/shared/logging"
)

const buildVersion string = "0.3.2"
const buildDate string = "2024-08-22"
const buildCommit string = "iter23"

func main() {
	err := logging.InitializeZapLogger("Info")
	if err != nil {
		log.Fatal(err)
	}
	logging.LogVersion(buildVersion, buildDate, buildCommit)
	serverFlags := new(config.ServerFlags)
	if err = config.ParseFlags(serverFlags); err != nil {
		log.Fatal(err)
	}
	err = app.RunServer(serverFlags)
	if err != nil {
		log.Fatal(err)
	}
}
