package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

type AgentFlags struct {
	Address        string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	HashKey        string `env:"KEY"`
	RateLimit      int    `env:"RATE_LIMIT"`
}

func ParseFlags(af *AgentFlags) {
	flag.StringVar(&af.Address, "a", "localhost:8080", "server address:port")
	flag.IntVar(&af.ReportInterval, "r", 10, "metric reporting interval, seconds int")
	flag.IntVar(&af.PollInterval, "p", 2, "metric polling interval, seconds int")
	flag.StringVar(&af.HashKey, "k", "", "SHA hash key")
	flag.IntVar(&af.RateLimit, "l", 1, "sending rate limit")

	flag.Parse()

	var ec AgentFlags
	if err := env.Parse(&ec); err != nil {
		panic(err)
	}
	if ec.Address != "" {
		af.Address = ec.Address
	}
	if ec.ReportInterval != 0 {
		af.ReportInterval = ec.ReportInterval
	}
	if ec.PollInterval != 0 {
		af.PollInterval = ec.PollInterval
	}
	if ec.HashKey != "" {
		af.HashKey = ec.HashKey
	}
	if ec.RateLimit != 0 {
		af.RateLimit = ec.RateLimit
	}
}
