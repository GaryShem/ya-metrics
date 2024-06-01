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

func ParseFlags(sf *ServerFlags) {
	flag.StringVar(&sf.Address, "a", "localhost:8080", "server address:port")
	flag.IntVar(&sf.StoreInterval, "i", 300, "metric saving interval")
	flag.StringVar(&sf.FileStoragePath, "f", "/tmp/metrics-db.json", "storage file path")
	flag.BoolVar(&sf.Restore, "r", true, "restore metrics from file")
	flag.StringVar(&sf.DBString, "d", "", "database connection string")
	flag.StringVar(&sf.HashKey, "k", "", "SHA hash key")
	flag.Parse()

	var ec ServerFlags
	if err := env.Parse(&ec); err != nil {
		panic(err)
	}
	if ec.Address != "" {
		sf.Address = ec.Address
	}
	if ec.FileStoragePath != "" {
		sf.FileStoragePath = ec.FileStoragePath
	}
	if ec.Restore {
		sf.Restore = ec.Restore
	}
	if ec.StoreInterval != 0 {
		sf.StoreInterval = ec.StoreInterval
	}
	if ec.DBString != "" {
		sf.DBString = ec.DBString
	}
	if ec.HashKey != "" {
		sf.HashKey = ec.HashKey
	}
}
