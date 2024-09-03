package config

import (
	"encoding/json"
	"flag"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
)

type ServerFlags struct {
	Address           string        `env:"ADDRESS" json:"address"`
	StoreInterval     int           `env:"STORE_INTERVAL" json:"-"`
	StoreIntervalJSON time.Duration `json:"store_interval"`
	FileStoragePath   string        `env:"FILE_STORAGE_PATH" json:"store_file"`
	Restore           bool          `env:"RESTORE" json:"restore"`
	DBString          string        `env:"DATABASE_DSN" json:"database_dsn"`
	HashKey           string        `env:"KEY" json:"hash_key"`
	CryptoKey         string        `env:"CRYPTO_KEY" json:"crypto_key"`
	Config            string        `env:"CONFIG"`
	TrustedSubnet     string        `env:"TRUSTED_SUBNET" json:"trusted_subnet"`
	GRPCAddress       string        `env:"GRPC_ADDRESS" json:"grpc_address"`
}

func withCmdLine() Option {
	return func(serverFlags *ServerFlags) error {
		flag.StringVar(&serverFlags.Address, "a", "localhost:8080", "server address:port")
		flag.IntVar(&serverFlags.StoreInterval, "i", 300, "metric saving interval")
		flag.StringVar(&serverFlags.FileStoragePath, "f", "/tmp/metrics-db.json", "storage file path")
		flag.BoolVar(&serverFlags.Restore, "r", true, "restore metrics from file")
		flag.StringVar(&serverFlags.DBString, "d", "", "database connection string")
		flag.StringVar(&serverFlags.HashKey, "k", "", "SHA hash key")
		flag.StringVar(&serverFlags.CryptoKey, "crypto-key", "", "crypto key")
		flag.StringVar(&serverFlags.Config, "c", "", "json config file")
		flag.StringVar(&serverFlags.TrustedSubnet, "t", "", "trusted subnet")
		flag.StringVar(&serverFlags.GRPCAddress, "g", "", "gRPC address")
		flag.Parse()
		return nil
	}
}

func withEnv() Option {
	return func(serverFlags *ServerFlags) error {
		var ec ServerFlags
		if err := env.Parse(&ec); err != nil {
			panic(err)
		}
		if ec.Address != "" {
			serverFlags.Address = ec.Address
		}
		if ec.FileStoragePath != "" {
			serverFlags.FileStoragePath = ec.FileStoragePath
		}
		if ec.Restore {
			serverFlags.Restore = ec.Restore
		}
		if ec.StoreInterval != 0 {
			serverFlags.StoreInterval = ec.StoreInterval
		}
		if ec.DBString != "" {
			serverFlags.DBString = ec.DBString
		}
		if ec.HashKey != "" {
			serverFlags.HashKey = ec.HashKey
		}
		if ec.CryptoKey != "" {
			serverFlags.CryptoKey = ec.CryptoKey
		}
		if ec.TrustedSubnet != "" {
			serverFlags.TrustedSubnet = ec.TrustedSubnet
		}
		if ec.GRPCAddress != "" {
			serverFlags.GRPCAddress = ec.GRPCAddress
		}
		return nil
	}
}

func withJSONConfig() Option {
	return func(serverFlags *ServerFlags) error {
		if serverFlags.Config == "" {
			return nil
		}
		config, err := os.ReadFile(serverFlags.Config)
		if err != nil {
			return err
		}
		var jsonFlags ServerFlags
		if err = json.Unmarshal(config, &jsonFlags); err != nil {
			return err
		}
		if serverFlags.Address == "" {
			serverFlags.Address = jsonFlags.Address
		}
		if serverFlags.StoreInterval == 0 {
			serverFlags.StoreInterval = int(jsonFlags.StoreIntervalJSON.Seconds())
		}
		if serverFlags.FileStoragePath == "" {
			serverFlags.FileStoragePath = jsonFlags.FileStoragePath
		}
		if !serverFlags.Restore {
			serverFlags.Restore = jsonFlags.Restore
		}
		if serverFlags.DBString == "" {
			serverFlags.DBString = jsonFlags.DBString
		}
		if serverFlags.HashKey == "" {
			serverFlags.HashKey = jsonFlags.HashKey
		}
		if serverFlags.CryptoKey == "" {
			serverFlags.CryptoKey = jsonFlags.CryptoKey
		}
		if serverFlags.TrustedSubnet == "" {
			serverFlags.TrustedSubnet = jsonFlags.TrustedSubnet
		}
		if serverFlags.GRPCAddress == "" {
			serverFlags.GRPCAddress = jsonFlags.GRPCAddress
		}
		return nil
	}
}

type Option func(*ServerFlags) error

func ParseFlags() (*ServerFlags, error) {
	sf := &ServerFlags{}
	options := []Option{withCmdLine(), withEnv(), withJSONConfig()}
	for _, option := range options {
		if err := option(sf); err != nil {
			return nil, err
		}
	}
	return sf, nil
}
