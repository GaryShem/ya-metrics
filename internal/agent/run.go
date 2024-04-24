package agent

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-resty/resty/v2"
)

type AgentFlags struct {
	Address        *string
	ReportInterval *int
	PollInterval   *int
}

func collectMetrics(mc *MetricCollector) error {
	return mc.CollectMetrics()
}

func sendMetrics(mc *MetricCollector, host string) error {
	client := resty.New()
	metrics, errDump := mc.DumpMetrics()
	if errDump != nil {
		return fmt.Errorf("error dumping metrics: %w", errDump)
	}
	url := "http://{host}/update"
	for _, m := range metrics {
		mJSON, err := json.Marshal(m)
		if err != nil {
			return fmt.Errorf("error marshalling metric: %w", err)
		}
		res, err := client.R().SetPathParam("host", host).
			SetHeader("Content-Type", "application/json").SetBody(mJSON).Post(url)
		if err != nil {
			return fmt.Errorf("error sending metric: %w", err)
		}
		if res.StatusCode() != 200 {
			return fmt.Errorf("error sending metric: %d %s", res.StatusCode(), res.String())
		}
	}
	return nil
}

func RunAgent(af *AgentFlags, sendOnce bool) error {
	metrics := NewMetricCollector(SupportedRuntimeMetrics())
	log.Printf("Server Address: %v\n", *af.Address)

	pollInterval := time.Second * time.Duration(*af.PollInterval)
	reportInterval := time.Second * time.Duration(*af.ReportInterval)

	collectionDelay := pollInterval
	dumpDelay := reportInterval
	log.Println("Starting metrics collection")
	for {
		sleepTime := min(dumpDelay, collectionDelay)
		log.Printf("Sleep %v\n", sleepTime)
		time.Sleep(sleepTime)
		dumpDelay -= sleepTime
		collectionDelay -= sleepTime
		if collectionDelay <= 0 {
			log.Println("collecting metrics")
			collectionDelay += pollInterval
			if err := collectMetrics(metrics); err != nil {
				return err
			}
		}
		if dumpDelay <= 0 {
			dumpDelay += reportInterval
			log.Println("sending metrics")
			if err := sendMetrics(metrics, *af.Address); err != nil {
				return err
			}
			if sendOnce {
				break
			}
		}
	}
	return nil
}
