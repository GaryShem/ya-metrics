package main

import (
	"log"

	"github.com/GaryShem/ya-metrics.git/internal/agent/app"
	"github.com/GaryShem/ya-metrics.git/internal/agent/config"
	"github.com/GaryShem/ya-metrics.git/internal/agent/metrics"
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
