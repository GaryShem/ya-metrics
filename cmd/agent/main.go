package main

import (
	"fmt"
	"github.com/GaryShem/ya-metrics.git/internal/storage"
	"net/http"
	"time"
)

var GaugeMetrics []string = []string{
	"Alloc",
	"BuckHashSys",
	"Frees",
	"GCCPUFraction",
	"GCSys",
	"HeapAlloc",
	"HeapIdle",
	"HeapInuse",
	"HeapObjects",
	"HeapReleased",
	"HeapSys",
	"LastGC",
	"Lookups",
	"MCacheInuse",
	"MCacheSys",
	"MSpanInuse",
	"MSpanSys",
	"Mallocs",
	"NextGC",
	"OtherSys",
	"PauseTotalNs",
	"StackInuse",
	"StackSys",
	"Sys",
	"TotalAlloc",
}

func collectMetrics(mc *storage.MetricCollector) {
	mc.CollectMetrics()
}

func sendMetrics(mc *storage.MetricCollector, host string) {
	metrics := mc.DumpMetrics()
	for name, value := range metrics.GaugeMetrics {
		requestURL := fmt.Sprintf("http://%v/update/gauge/%v/%v", host, name, value)
		res, err := http.Post(requestURL, "text/plain", nil)
		if err != nil {
			panic(err)
		}
		if res.StatusCode != 200 {
			panic(res.StatusCode)
		}
		if err := res.Body.Close(); err != nil {
			panic(err)
		}
	}
	for name, value := range metrics.CounterMetrics {
		requestURL := fmt.Sprintf("http://%v/update/counter/%v/%v", host, name, value)
		res, err := http.Post(requestURL, "text/plain", nil)
		if err != nil {
			panic(err)
		}
		if res.StatusCode != 200 {
			panic(res.StatusCode)
		}
		if err := res.Body.Close(); err != nil {
			panic(err)
		}
	}
}

func main() {
	metrics := storage.NewMetricCollector(GaugeMetrics)

	collectionPeriod := time.Second * 2
	dumpPeriod := time.Second * 10

	collectionDelay := collectionPeriod
	dumpDelay := dumpPeriod

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
			collectionDelay += collectionPeriod
			collectMetrics(metrics)
		}
		if dumpDelay <= 0 {
			dumpDelay += dumpPeriod
			fmt.Println("sending metrics")
			sendMetrics(metrics, `localhost:8080`)
		}
	}
}
