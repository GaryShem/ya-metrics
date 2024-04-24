package main

import (
	"github.com/GaryShem/ya-metrics.git/internal/server"
	"github.com/GaryShem/ya-metrics.git/internal/shared/storage/memStorage"
)

func main() {
	ms := memStorage.NewMemStorage()
	sf := new(server.ServerFlags)
	ParseFlags(sf)
	server.RunServer(sf, ms)
}
