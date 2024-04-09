package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/GaryShem/ya-metrics.git/internal/storage"
)

func collectMetrics(mc *storage.MetricCollector) {
	mc.CollectMetrics()
}

func sendMetrics(mc *storage.MetricCollector, host string) {
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
			panic(err)
		}
		if res.StatusCode() != 200 {
			panic(res.StatusCode)
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
			panic(err)
		}
		if res.StatusCode() != 200 {
			panic(res.StatusCode)
		}
	}
}

func main() {
	metrics := storage.NewMetricCollector(storage.RuntimeMetrics)
	af := new(AgentFlags)
	ParseFlags(af)

	pollInterval := time.Second * time.Duration(*af.pollInterval)
	reportInterval := time.Second * time.Duration(*af.reportInterval)

	collectionDelay := pollInterval
	dumpDelay := reportInterval

	fmt.Println("Starting metrics collection")
	var i int64 = 0
	for {
		fmt.Printf("Iteration %v\n", i)
		i += 1
		sleepTime := min(dumpDelay, collectionDelay)
		fmt.Printf("Sleep %v\n", sleepTime)
		time.Sleep(sleepTime)
		dumpDelay -= sleepTime
		collectionDelay -= sleepTime
		if collectionDelay <= 0 {
			fmt.Println("collecting metrics")
			collectionDelay += pollInterval
			collectMetrics(metrics)
		}
		if dumpDelay <= 0 {
			dumpDelay += reportInterval
			fmt.Println("sending metrics")
			sendMetrics(metrics, *af.address)
		}
	}
}
