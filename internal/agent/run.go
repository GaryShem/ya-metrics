package agent

import (
	"bytes"
	gzip "compress/gzip"
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

func triggerCollectMetrics(mc *MetricCollector, interval time.Duration, ec chan error) {
	timer := time.NewTicker(interval)
	defer timer.Stop()
	for {
		<-timer.C
		err := mc.CollectMetrics()
		if err != nil {
			ec <- err
			return
		}
	}

}

func triggerSendMetrics(mc *MetricCollector, host string, sendOnce bool, ignoreSendError bool,
	gzipRequest bool, interval time.Duration, ec chan error) {
	timer := time.NewTicker(interval)
	defer timer.Stop()
	for {
		<-timer.C
		err := sendMetrics(mc, host, gzipRequest)
		if err != nil {
			if ignoreSendError {
				logging.Log.Warnln("Ignore send error:", err)
				continue
			}
			ec <- err
			return
		}
		if sendOnce {
			ec <- nil
		}
	}
}

func wrapGzipRequest(r *resty.Request, mJSON []byte) error {
	//copy(bodyCopy, mJSON)
	var buffer bytes.Buffer
	writer, err := gzip.NewWriterLevel(&buffer, gzip.BestCompression)
	if err != nil {
		return fmt.Errorf("failed init compress writer: %v", err)
	}
	if _, err := writer.Write(mJSON); err != nil {
		return fmt.Errorf("error gzipping metric: %w", err)
	}
	if err := writer.Close(); err != nil {
		return fmt.Errorf("error gzipping metric: %w", err)
	}
	bodyBytes := buffer.Bytes()
	r.SetBody(bodyBytes)
	r.Header.Add("Content-Encoding", "gzip")
	r.Header.Add("Accept-Encoding", "gzip")
	return nil
}

func sendMetrics(mc *MetricCollector, host string, gzipRequest bool) error {
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
			SetHeader("Content-Type", "application/json")
		if gzipRequest {
			err = wrapGzipRequest(request, mJSON)
			if err != nil {
				return fmt.Errorf("error gzipping metric: %w", err)
			}
		} else {
			request.SetBody(mJSON)
		}
		res, err := request.Post(url)
		if err != nil {
			if res != nil {
				return fmt.Errorf("error sending metric, response is not nil: %w, %d %s", err, res.StatusCode(), res.String())
			} else {
				return fmt.Errorf("error sending metric, response is nil: %w", err)
			}
		}
		if res.StatusCode() != http.StatusOK {
			return fmt.Errorf("status code not 200: %d %s", res.StatusCode(), res.String())
		}
	}
	return nil
}

func RunAgent(af *AgentFlags, runtimeMetrics []string, sendOnce bool, ignoreSendError bool, gzipRequest bool) error {
	if err := logging.InitializeZapLogger("Info"); err != nil {
		return fmt.Errorf("error initializing logger: %w", err)
	}
	logging.Log.Infoln("agent started")
	metrics := NewMetricCollector(runtimeMetrics)
	logging.Log.Infoln("Server Address:", *af.Address)

	pollInterval := time.Second * time.Duration(*af.PollInterval)
	reportInterval := time.Second * time.Duration(*af.ReportInterval)

	log.Println("Starting metrics collection")
	c := make(chan error)
	go triggerCollectMetrics(metrics, pollInterval, c)
	go triggerSendMetrics(metrics, *af.Address, sendOnce, ignoreSendError, gzipRequest, reportInterval, c)
	return <-c
}
