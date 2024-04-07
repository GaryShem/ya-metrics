package main

import (
	"flag"

	"github.com/caarlos0/env/v6"

	"github.com/GaryShem/ya-metrics.git/internal/server"
)

type envConfig struct {
	Address string `env:"ADDRESS"`
}

func ParseFlags(sf *server.ServerFlags) {
	sf.Address = flag.String("a", "localhost:8080", "server address:port")
	flag.Parse()

	ec := envConfig{}
	if err := env.Parse(&ec); err != nil {
		panic(err)
	}
	if ec.Address != "" {
		sf.Address = &ec.Address
	}
}
