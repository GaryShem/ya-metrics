package agent

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
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
	metrics := mc.DumpMetrics()
	url := "http://{host}/update/{type}/{name}/{value}"
	for name, value := range metrics.GaugeMetrics {
		request := client.R().SetPathParams(map[string]string{
			"host":  host,
			"type":  "gauge",
			"name":  name,
			"value": strconv.FormatFloat(value, 'f', 6, 64),
		})
		res, err := request.Post(url)
		if err != nil {
			return fmt.Errorf("could not send metrics: %w", err)
		}
		if res.StatusCode() != http.StatusOK {
			return fmt.Errorf("could not send metrics, return code: %v",
				res.StatusCode())
		}
	}
	for name, value := range metrics.CounterMetrics {
		request := client.R().SetPathParams(map[string]string{
			"host":  host,
			"type":  "counter",
			"name":  name,
			"value": strconv.FormatInt(value, 10),
		})
		res, err := request.Post(url)
		if err != nil {
			return fmt.Errorf("could not send metrics: %w", err)
		}
		if res.StatusCode() != http.StatusOK {
			return fmt.Errorf("could not send metrics, return code: %v",
				res.StatusCode())
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
