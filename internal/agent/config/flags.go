package config

import (
	"encoding/json"
	"flag"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
)

type AgentFlags struct {
	Address            string        `env:"ADDRESS" json:"address"`
	ReportInterval     int           `env:"REPORT_INTERVAL"`
	PollInterval       int           `env:"POLL_INTERVAL"`
	ReportIntervalJSON time.Duration `json:"report_interval"`
	PollIntervalJSON   time.Duration `json:"poll_interval"`
	HashKey            string        `env:"KEY" json:"hash_key"`
	RateLimit          int           `env:"RATE_LIMIT" json:"rate_limit"`
	GzipRequest        bool          `env:"GZIP_REQUEST" json:"gzip_request"`
	CryptoKey          string        `env:"CRYPTO_KEY" json:"crypto_key"`
	Config             string        `json:"config"`
	GRPCAddress        string        `env:"GRPC_ADDRESS" json:"grpc_address"`
}

func withCmdLine() Option {
	return func(agentFlags *AgentFlags) error {
		flag.StringVar(&agentFlags.Address, "a", "localhost:8080", "server address:port")
		flag.IntVar(&agentFlags.ReportInterval, "r", 10, "metric reporting interval, seconds int")
		flag.IntVar(&agentFlags.PollInterval, "p", 2, "metric polling interval, seconds int")
		flag.StringVar(&agentFlags.HashKey, "k", "", "SHA hash key")
		flag.IntVar(&agentFlags.RateLimit, "l", 1, "sending rate limit")
		flag.BoolVar(&agentFlags.GzipRequest, "z", true, "gzip request")
		flag.StringVar(&agentFlags.CryptoKey, "crypto-key", "", "crypto key")
		flag.StringVar(&agentFlags.Config, "c", "", "json config file")
		flag.StringVar(&agentFlags.GRPCAddress, "g", "", "gRPC address")
		flag.Parse()
		return nil
	}
}

func withEnv() Option {
	return func(agentFlags *AgentFlags) error {
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
		if ec.GRPCAddress != "" {
			agentFlags.GRPCAddress = ec.GRPCAddress
		}
		return nil
	}
}

func withJSONConfig() Option {
	return func(agentFlags *AgentFlags) error {
		if agentFlags.Config == "" {
			return nil
		}
		config, err := os.ReadFile(agentFlags.Config)
		if err != nil {
			return err
		}
		var jsonFlags AgentFlags
		if err = json.Unmarshal(config, &jsonFlags); err != nil {
			return err
		}
		if agentFlags.Address == "" {
			agentFlags.Address = jsonFlags.Address
		}
		if agentFlags.PollInterval == 0 {
			agentFlags.PollInterval = int(jsonFlags.PollIntervalJSON.Seconds())
		}
		if agentFlags.ReportInterval == 0 {
			agentFlags.ReportInterval = int(jsonFlags.ReportIntervalJSON.Seconds())
		}
		if !agentFlags.GzipRequest {
			agentFlags.GzipRequest = jsonFlags.GzipRequest
		}
		if agentFlags.RateLimit == 0 {
			agentFlags.RateLimit = jsonFlags.RateLimit
		}
		if agentFlags.HashKey == "" {
			agentFlags.HashKey = jsonFlags.HashKey
		}
		if agentFlags.CryptoKey == "" {
			agentFlags.CryptoKey = jsonFlags.CryptoKey
		}
		if agentFlags.GRPCAddress == "" {
			agentFlags.GRPCAddress = jsonFlags.GRPCAddress
		}
		return nil
	}
}

type Option func(*AgentFlags) error

func ParseFlags() (*AgentFlags, error) {
	af := &AgentFlags{}
	options := []Option{withCmdLine(), withEnv(), withJSONConfig()}
	for _, option := range options {
		if err := option(af); err != nil {
			return nil, err
		}
	}
	return af, nil
}
