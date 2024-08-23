package main

import (
	"log"

	"github.com/GaryShem/ya-metrics.git/internal/server/app"
	"github.com/GaryShem/ya-metrics.git/internal/server/config"
	"github.com/GaryShem/ya-metrics.git/internal/shared/logging"
)

const buildVersion string = "N/A"
const buildDate string = "N/A"
const buildCommit string = "N/A"

func main() {
	err := logging.InitializeZapLogger("Info")
	if err != nil {
		log.Fatal(err)
	}
	logging.LogVersion(buildVersion, buildDate, buildCommit)
	serverFlags, err := config.ParseFlags()
	if err != nil {
		log.Fatal(err)
	}
	err = app.RunServer(serverFlags)
	if err != nil {
		log.Fatal(err)
	}
}
