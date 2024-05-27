package main

import (
	"flag"

	"github.com/caarlos0/env/v6"

	"github.com/GaryShem/ya-metrics.git/internal/server"
)

type envConfig struct {
	Address         string `env:"ADDRESS"`
	StoreInterval   *int   `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         *bool  `env:"RESTORE"`
	DBString        string `env:"DATABASE_DSN"`
}

func ParseFlags(sf *server.ServerFlags) {
	sf.Address = flag.String("a", "localhost:8080", "server address:port")
	sf.StoreInterval = flag.Int("i", 300, "metric saving interval")
	sf.FileStoragePath = flag.String("f", "/tmp/metrics-db.json", "storage file path")
	sf.Restore = flag.Bool("r", true, "restore metrics from file")
	sf.DBString = flag.String("d", "", "database connection string")
	flag.Parse()

	ec := envConfig{}
	if err := env.Parse(&ec); err != nil {
		panic(err)
	}
	if ec.Address != "" {
		sf.Address = &ec.Address
	}
	if ec.FileStoragePath != "" {
		sf.FileStoragePath = &ec.FileStoragePath
	}
	if ec.Restore != nil {
		sf.Restore = ec.Restore
	}
	if ec.StoreInterval != nil {
		sf.StoreInterval = ec.StoreInterval
	}
	if ec.DBString != "" {
		sf.DBString = &ec.DBString
	}
}
