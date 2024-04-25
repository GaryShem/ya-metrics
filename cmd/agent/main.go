package main

import (
	"log"

	"github.com/GaryShem/ya-metrics.git/internal/agent"
	"github.com/GaryShem/ya-metrics.git/internal/shared/logging"
)

func main() {
	err := logging.InitializeZapLogger("Info")
	if err != nil {
		log.Fatal(err)
	}
	af := new(agent.AgentFlags)
	ParseFlags(af)
	logging.Log.Infoln("client starting with flags",
		"host", *af.Address,
		"poll interval", *af.PollInterval,
		"send interval", *af.ReportInterval,
	)
	err = agent.RunAgent(af, false)
	if err != nil {
		log.Fatalf("agent closed with error %v", err)
	}
}
