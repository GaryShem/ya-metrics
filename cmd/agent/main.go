package main

import (
	"log"

	"github.com/GaryShem/ya-metrics.git/internal/agent/app"
	"github.com/GaryShem/ya-metrics.git/internal/agent/config"
	"github.com/GaryShem/ya-metrics.git/internal/agent/metrics"
	"github.com/GaryShem/ya-metrics.git/internal/shared/logging"
)

const buildVersion string = "0.20"
const buildDate string = "2024-08-07"
const buildCommit string = ""

func main() {
	err := logging.InitializeZapLogger("Info")
	if err != nil {
		log.Fatal(err)
	}
	logging.LogVersion(buildVersion, buildDate, buildCommit)
	agentFlags := new(config.AgentFlags)
	config.ParseFlags(agentFlags)
	logging.Log.Infoln("client starting with flags",
		"host", agentFlags.Address,
		"poll interval", agentFlags.PollInterval,
		"send interval", agentFlags.ReportInterval,
		"hash key", agentFlags.HashKey,
	)
	err = app.RunAgent(agentFlags, metrics.SupportedRuntimeMetrics(), false, true)
	if err != nil {
		log.Fatalf("agent closed with error %v", err)
	}
}
