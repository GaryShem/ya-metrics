package main

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

type envConfig struct {
	Address        string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
}

type AgentFlags struct {
	address        *string
	reportInterval *int
	pollInterval   *int
}

func ParseFlags(af *AgentFlags) {
	af.address = flag.String("a", "localhost:8080", "server address:port")
	af.reportInterval = flag.Int("r", 10, "metric reporting interval")
	af.pollInterval = flag.Int("p", 2, "metric polling interval")
	flag.Parse()

	var ec envConfig
	if err := env.Parse(&ec); err != nil {
		panic(err)
	}
	if ec.Address != "" {
		af.address = &ec.Address
	}
	if ec.ReportInterval != 0 {
		af.reportInterval = &ec.ReportInterval
	}
	if ec.PollInterval != 0 {
		af.pollInterval = &ec.PollInterval
	}
}
