package main

import "flag"

type AgentFlags struct {
	address        string
	reportInterval int
	pollInterval   int
}

func ParseFlags(af *AgentFlags) {
	af.address = *flag.String("a", "localhost:8080", "server address:port")
	af.reportInterval = *flag.Int("r", 10, "metric reporting interval")
	af.pollInterval = *flag.Int("p", 2, "metric polling interval")
	flag.Parse()
}
