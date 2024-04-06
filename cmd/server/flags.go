package main

import "flag"

type ServerFlags struct {
	address string
}

func ParseFlags(sf *ServerFlags) {
	sf.address = *flag.String("a", "localhost:8080", "server address:port")
	flag.Parse()
}
