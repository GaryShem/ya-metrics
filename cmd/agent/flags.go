package main

import (
	"flag"
	"os"
	"strconv"
)

//type envConfig struct {
//	address        string `env:"ADDRESS"`
//	reportInterval int    `env:"REPORT_INTERVAL"`
//	pollInterval   int    `env:"POLL_INTERVAL"`
//}

type AgentFlags struct {
	address        *string
	reportInterval *int
	pollInterval   *int
}

func ParseFlags(af *AgentFlags) {
	af.address = flag.String("a", "localhost:8080", "server address:port")
	af.reportInterval = flag.Int("r", 10, "metric reporting interval")
	af.pollInterval = flag.Int("p", 2, "metric polling interval")
	flag.Parse()

	addressEnv := os.Getenv("ADDRESS")
	if addressEnv != "" {
		af.address = &addressEnv
	}
	reportIntervalEnv := os.Getenv("REPORT_INTERVAL")
	if reportIntervalEnv != "" {
		reportInterval, err := strconv.Atoi(reportIntervalEnv)
		if err == nil {
			af.reportInterval = &reportInterval
		}
	}
	pollIntervalEnv := os.Getenv("POLL_INTERVAL")
	if pollIntervalEnv != "" {
		pollInterval, err := strconv.Atoi(reportIntervalEnv)
		if err == nil {
			af.pollInterval = &pollInterval
		}
	}

	//var ec envConfig
	//if err := env.Parse(&ec); err != nil {
	//	panic(err)
	//}
	//
	//if ec.address != "" {
	//	af.address = &ec.address
	//}
	//if ec.reportInterval != 0 {
	//	af.reportInterval = &ec.reportInterval
	//}
	//if ec.pollInterval != 0 {
	//	af.pollInterval = &ec.pollInterval
	//}
}
