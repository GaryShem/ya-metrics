package main

import (
	"log"

	"github.com/GaryShem/ya-metrics.git/internal/server/app"
	"github.com/GaryShem/ya-metrics.git/internal/server/config"
	"github.com/GaryShem/ya-metrics.git/internal/shared/logging"
)

const buildVersion string = "0.3.0"
const buildDate string = "2024-08-21"
const buildCommit string = "iter21"

func main() {
	err := logging.InitializeZapLogger("Info")
	if err != nil {
		log.Fatal(err)
	}
	logging.LogVersion(buildVersion, buildDate, buildCommit)
	serverFlags := new(config.ServerFlags)
	config.ParseFlags(serverFlags)
	err = app.RunServer(serverFlags)
	if err != nil {
		log.Fatal(err)
	}
}
