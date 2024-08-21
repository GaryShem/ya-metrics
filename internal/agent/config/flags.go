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
	GzipRequest    bool
	CryptoKey      string `env:"CRYPTO_KEY"`
}

func ParseFlags(agentFlags *AgentFlags) {
	flag.StringVar(&agentFlags.Address, "a", "localhost:8080", "server address:port")
	flag.IntVar(&agentFlags.ReportInterval, "r", 10, "metric reporting interval, seconds int")
	flag.IntVar(&agentFlags.PollInterval, "p", 2, "metric polling interval, seconds int")
	flag.StringVar(&agentFlags.HashKey, "k", "", "SHA hash key")
	flag.IntVar(&agentFlags.RateLimit, "l", 1, "sending rate limit")
	flag.BoolVar(&agentFlags.GzipRequest, "z", true, "gzip request")
	flag.StringVar(&agentFlags.CryptoKey, "crypto-key", "", "crypto key")

	flag.Parse()

	var ec AgentFlags
	if err := env.Parse(&ec); err != nil {
		panic(err)
	}
	if ec.Address != "" {
		agentFlags.Address = ec.Address
	}
	if ec.ReportInterval != 0 {
		agentFlags.ReportInterval = ec.ReportInterval
	}
	if ec.PollInterval != 0 {
		agentFlags.PollInterval = ec.PollInterval
	}
	if ec.HashKey != "" {
		agentFlags.HashKey = ec.HashKey
	}
	if ec.RateLimit != 0 {
		agentFlags.RateLimit = ec.RateLimit
	}
	if ec.CryptoKey != "" {
		agentFlags.CryptoKey = ec.CryptoKey
	}
}
