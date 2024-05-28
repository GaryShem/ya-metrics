package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

type ServerFlags struct {
	Address         string
	StoreInterval   int
	FileStoragePath string
	Restore         bool
	DBString        string
	HashKey         string
}

type envConfig struct {
	Address         string `env:"ADDRESS"`
	StoreInterval   *int   `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         *bool  `env:"RESTORE"`
	DBString        string `env:"DATABASE_DSN"`
	HashKey         string `env:"KEY"`
}

type flagConfig struct {
	Address         *string
	StoreInterval   *int
	FileStoragePath *string
	Restore         *bool
	DBString        *string
	HashKey         *string
}

func ParseFlags(sf *ServerFlags) {
	flags := &flagConfig{
		Address:         flag.String("a", "localhost:8080", "server address:port"),
		StoreInterval:   flag.Int("i", 300, "metric saving interval"),
		FileStoragePath: flag.String("f", "/tmp/metrics-db.json", "storage file path"),
		Restore:         flag.Bool("r", true, "restore metrics from file"),
		DBString:        flag.String("d", "", "database connection string"),
		HashKey:         flag.String("k", "", "SHA hash key"),
	}
	flag.Parse()

	ec := envConfig{}
	if err := env.Parse(&ec); err != nil {
		panic(err)
	}
	if ec.Address != "" {
		sf.Address = ec.Address
	} else {
		sf.Address = *flags.Address
	}
	if ec.FileStoragePath != "" {
		sf.FileStoragePath = ec.FileStoragePath
	} else {
		sf.FileStoragePath = *flags.FileStoragePath
	}
	if ec.Restore != nil {
		sf.Restore = *ec.Restore
	} else {
		sf.Restore = *flags.Restore
	}
	if ec.StoreInterval != nil {
		sf.StoreInterval = *ec.StoreInterval
	} else {
		sf.StoreInterval = *flags.StoreInterval
	}
	if ec.DBString != "" {
		sf.DBString = ec.DBString
	} else {
		sf.DBString = *flags.DBString
	}
	if ec.HashKey != "" {
		sf.HashKey = ec.HashKey
	} else {
		sf.HashKey = *flags.HashKey
	}
}
