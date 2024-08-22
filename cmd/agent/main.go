package main

import (
	"log"

	"github.com/GaryShem/ya-metrics.git/internal/agent/app"
	"github.com/GaryShem/ya-metrics.git/internal/agent/config"
	"github.com/GaryShem/ya-metrics.git/internal/agent/metrics"
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
	agentFlags, err := config.ParseFlags()
	if err != nil {
		log.Fatal(err)
	}
	logging.Log.Infoln("client starting with flags",
		"host", agentFlags.Address,
		"poll interval", agentFlags.PollInterval,
		"send interval", agentFlags.ReportInterval,
		"hash key", agentFlags.HashKey,
	)
	err = app.RunAgent(agentFlags, metrics.SupportedRuntimeMetrics(), false)
	if err != nil {
		log.Fatalf("agent closed with error %v", err)
	}
}
