package main

import (
	"log"

	"github.com/GaryShem/ya-metrics.git/internal/server"
)

func main() {
	sf := new(server.ServerFlags)
	ParseFlags(sf)
	err := server.RunServer(sf)
	if err != nil {
		log.Fatal(err)
	}
}
