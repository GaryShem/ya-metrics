package agent

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/GaryShem/ya-metrics.git/internal/shared/logging"
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
	logging.Log.Infoln(host)
	if errDump != nil {
		return fmt.Errorf("error dumping metrics: %w", errDump)
	}
	url := "http://{host}/update"

	for _, m := range metrics {
		mJSON, err := json.Marshal(m)
		if err != nil {
			return fmt.Errorf("error marshalling metric: %w", err)
		}
		request := client.R().SetPathParam("host", host).
			SetHeader("Content-Type", "application/json").SetBody(mJSON)
		res, err := request.Post(url)
		if err != nil {
			return fmt.Errorf("error sending metric: %w", err)
		}
		if res.StatusCode() != http.StatusOK {
			return fmt.Errorf("error sending metric: %d %s", res.StatusCode(), res.String())
		}
	}
	return nil
}

func RunAgent(af *AgentFlags, sendOnce bool) error {
	if err := logging.InitializeZapLogger("Info"); err != nil {
		return fmt.Errorf("error initializing logger: %w", err)
	}
	logging.Log.Infoln("agent started")
	//log.Fatal("crash agent for science")
	metrics := NewMetricCollector(SupportedRuntimeMetrics())
	logging.Log.Infoln("Server Address:", *af.Address)

	pollInterval := time.Second * time.Duration(*af.PollInterval)
	reportInterval := time.Second * time.Duration(*af.ReportInterval)

	collectionDelay := pollInterval
	dumpDelay := reportInterval
	log.Println("Starting metrics collection")
	for {
		sleepTime := min(dumpDelay, collectionDelay)
		logging.Log.Infoln("Sleep", sleepTime)
		time.Sleep(sleepTime)
		dumpDelay -= sleepTime
		collectionDelay -= sleepTime
		if collectionDelay <= 0 {
			logging.Log.Infoln("collecting metrics")
			collectionDelay += pollInterval
			if err := collectMetrics(metrics); err != nil {
				return err
			}
		}
		if dumpDelay <= 0 {
			dumpDelay += reportInterval
			logging.Log.Infoln("sending metrics")
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
