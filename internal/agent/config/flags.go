package config

import (
	"flag"

	"github.com/caarlos0/env/v6"

	"github.com/GaryShem/ya-metrics.git/internal/agent/app"
)

type envConfig struct {
	Address        string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	HashKey        string `env:"KEY"`
}

func ParseFlags(af *app.AgentFlags) {
	af.Address = flag.String("a", "localhost:8080", "server address:port")
	af.ReportInterval = flag.Int("r", 10, "metric reporting interval, seconds int")
	af.PollInterval = flag.Int("p", 2, "metric polling interval, seconds int")
	af.HashKey = flag.String("k", "", "SHA hash key")

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
	if ec.HashKey != "" {
		af.HashKey = &ec.HashKey
	}
}
