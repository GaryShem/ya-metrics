package main

import (
	"flag"

	"github.com/caarlos0/env/v6"

	"github.com/GaryShem/ya-metrics.git/internal/agent"
)

type envConfig struct {
	Address        string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
}

func ParseFlags(af *agent.AgentFlags) {
	af.Address = flag.String("a", "localhost:8080", "server address:port")
	af.ReportInterval = flag.Int("r", 2, "metric reporting interval, seconds int")
	af.PollInterval = flag.Int("p", 1, "metric polling interval, seconds int")
	flag.Parse()

	var ec envConfig
	if err := env.Parse(&ec); err != nil {
		panic(err)
	}
	if ec.Address != "" {
		af.Address = &ec.Address
	}
	if ec.ReportInterval != 0 {
		af.ReportInterval = &ec.ReportInterval
	}
	if ec.PollInterval != 0 {
		af.PollInterval = &ec.PollInterval
	}
}
