package main

import (
	"log"

	"github.com/GaryShem/ya-metrics.git/internal/server/app"
	"github.com/GaryShem/ya-metrics.git/internal/server/config"
	"github.com/GaryShem/ya-metrics.git/internal/shared/logging"
)

var buildVersion string = "0.20"
var buildDate string = "2024-08-07"
var buildCommit string

func logVersion() {
	output := "N/A"
	if buildVersion != "" {
		output = buildVersion
	}
	logging.Log.Infof("Build version: %s", output)
	output = "N/A"
	if buildDate != "N/A" {
		output = buildDate
	}
	logging.Log.Infof("Build date: %s", output)
	output = "N/A"
	if buildCommit != "" {
		output = buildCommit
	}
	logging.Log.Infof("Build commit: %s", output)
}

func main() {
	err := logging.InitializeZapLogger("Info")
	if err != nil {
		log.Fatal(err)
	}
	logVersion()
	serverFlags := new(config.ServerFlags)
	config.ParseFlags(serverFlags)
	err = app.RunServer(serverFlags)
	if err != nil {
		log.Fatal(err)
	}
}
