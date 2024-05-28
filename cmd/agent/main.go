package main

import (
	"log"

	"github.com/GaryShem/ya-metrics.git/internal/agent/app"
	"github.com/GaryShem/ya-metrics.git/internal/agent/config"
	"github.com/GaryShem/ya-metrics.git/internal/agent/metrics"
	"github.com/GaryShem/ya-metrics.git/internal/shared/logging"
)

func main() {
	err := logging.InitializeZapLogger("Info")
	if err != nil {
		log.Fatal(err)
	}
	af := new(config.AgentFlags)
	config.ParseFlags(af)
	logging.Log.Infoln("client starting with flags",
		"host", *af.Address,
		"poll interval", *af.PollInterval,
		"send interval", *af.ReportInterval,
		"hash key", *af.HashKey,
	)
	err = app.RunAgent(af, metrics.SupportedRuntimeMetrics(),
		false, false, true)
	if err != nil {
		log.Fatalf("agent closed with error %v", err)
	}
}
