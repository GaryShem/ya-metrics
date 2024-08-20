package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

type ServerFlags struct {
	Address         string `env:"ADDRESS"`
	StoreInterval   int    `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
	DBString        string `env:"DATABASE_DSN"`
	HashKey         string `env:"KEY"`
}

func ParseFlags(serverFlags *ServerFlags) {
	flag.StringVar(&serverFlags.Address, "a", "localhost:8080", "server address:port")
	flag.IntVar(&serverFlags.StoreInterval, "i", 300, "metric saving interval")
	flag.StringVar(&serverFlags.FileStoragePath, "f", "/tmp/metrics-db.json", "storage file path")
	flag.BoolVar(&serverFlags.Restore, "r", true, "restore metrics from file")
	flag.StringVar(&serverFlags.DBString, "d", "", "database connection string")
	flag.StringVar(&serverFlags.HashKey, "k", "", "SHA hash key")
	flag.Parse()

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
}
