package main

import (
	"github.com/GaryShem/ya-metrics.git/internal/server"
)

func main() {
	ms := mem_storage.NewMemStorage()
	sf := new(server.ServerFlags)
	ParseFlags(sf)
	server.RunServer(sf, ms)
}
