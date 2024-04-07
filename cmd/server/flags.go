package main

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

type ServerFlags struct {
	address *string
}

type envConfig struct {
	Address string `env:"ADDRESS"`
}

func ParseFlags(sf *ServerFlags) {
	sf.address = flag.String("a", "localhost:8080", "server address:port")
	flag.Parse()

	ec := envConfig{}
	if err := env.Parse(&ec); err != nil {
		panic(err)
	}
	if ec.Address != "" {
		sf.address = &ec.Address
	}
}
