package main

import (
	"github.com/GaryShem/ya-metrics.git/internal/server"
)

func main() {
	sf := new(server.ServerFlags)
	ParseFlags(sf)
	server.RunServer(sf)
}
