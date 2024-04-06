package main

import (
	"flag"
	"fmt"
	"os"
)

type ServerFlags struct {
	address *string
}

//type envConfig struct {
//	address string `env:"ADDRESS"`
//}

func ParseFlags(sf *ServerFlags) {
	sf.address = flag.String("a", "localhost:8080", "server address:port")
	flag.Parse()

	addressEnv := os.Getenv("ADDRESS")
	fmt.Printf("ADDRESS env: %v\n", os.Getenv("ADDRESS"))

	if addressEnv != "" {
		sf.address = &addressEnv
	}

	//ec := envConfig{}
	//if err := env.Parse(&ec); err != nil {
	//	panic(err)
	//}
	//if ec.address != "" {
	//	sf.address = &ec.address
	//}
}
