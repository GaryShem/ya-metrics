package main

import (
	"log"

	"github.com/GaryShem/ya-metrics.git/internal/agent"
)

func main() {
	af := new(agent.AgentFlags)
	ParseFlags(af)
	err := agent.RunAgent(af, false)
	if err != nil {
		log.Fatal("agent closed with error")
	}
}
