package main

import (
	"github.com/GaryShem/ya-metrics.git/internal/server"
	"github.com/GaryShem/ya-metrics.git/internal/shared/storage/memstorage"
)

func main() {
	ms := memstorage.NewMemStorage()
	sf := new(server.ServerFlags)
	ParseFlags(sf)
	server.RunServer(sf, ms)
}
