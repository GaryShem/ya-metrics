package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/caarlos0/env/v6"
)

type ServerFlags struct {
	address *string
}

type envConfig struct {
	address string `env:"ADDRESS,required"`
}

func ParseFlags(sf *ServerFlags) {
	sf.address = flag.String("a", "localhost:8080", "server address:port")
	flag.Parse()

	fmt.Printf("ADDRESS env: %v\n", os.Getenv("ADDRESS"))
	ec := envConfig{}
	if err := env.Parse(&ec); err != nil {
		panic(err)
	}

	if ec.address != "" {
		sf.address = &ec.address
	}
}
