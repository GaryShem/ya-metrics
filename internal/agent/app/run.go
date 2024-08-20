package app

import (
	"bytes"
	gzip "compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/GaryShem/ya-metrics.git/internal/agent/config"
	"github.com/GaryShem/ya-metrics.git/internal/agent/metrics"
	"github.com/GaryShem/ya-metrics.git/internal/shared/logging"
	"github.com/GaryShem/ya-metrics.git/internal/shared/storage/models"
)

func CollectMetrics(mc *metrics.MetricCollector, interval time.Duration, ec chan error) {
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

func CollectAdditionalMetrics(mc *metrics.MetricCollector, interval time.Duration, ec chan error) {
	timer := time.NewTicker(interval)
	defer timer.Stop()
	for {
		<-timer.C
		err := mc.CollectAdditionalMetrics()
		if err != nil {
			ec <- err
			return
		}
	}
}

func SendMetrics(mc *metrics.MetricCollector, agentFlags config.AgentFlags, sendOnce bool, ignoreSendError bool, ec chan error) {
	defer logging.Log.Infoln("stopping sending metrics")
	interval := time.Duration(agentFlags.ReportInterval) * time.Second
	timer := time.NewTicker(interval)
	semaphore := make(chan struct{}, agentFlags.RateLimit)
	sendErrChan := make(chan error)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			semaphore <- struct{}{}
			// if we only need to send a single message (i.e. for tests), fill the buffer channel
			// that way we can ensure the sending goroutine is done with its task
			if sendOnce {
				for range agentFlags.RateLimit - 1 {
					semaphore <- struct{}{}
				}
			}
			metricsDump, err := mc.DumpMetrics()
			if err != nil {
				ec <- fmt.Errorf("error dumping metrics: %w", err)
				return
			}
			go sendMetricsBatch(metricsDump, agentFlags.Address, agentFlags.GzipRequest, agentFlags.HashKey, sendErrChan, semaphore)
			if sendOnce {
				semaphore <- struct{}{}
				ec <- nil
				return
			}
		case err := <-sendErrChan:
			ec <- err
			return
		}
	}
}

func wrapGzipBody(mJSON []byte) ([]byte, error) {
	var buffer bytes.Buffer
	writer, err := gzip.NewWriterLevel(&buffer, gzip.BestCompression)
	if err != nil {
		return nil, fmt.Errorf("failed init compress writer: %v", err)
	}
	if _, err = writer.Write(mJSON); err != nil {
		return nil, fmt.Errorf("error gzipping metric: %w", err)
	}
	if err = writer.Close(); err != nil {
		return nil, fmt.Errorf("error gzipping metric: %w", err)
	}
	bodyBytes := buffer.Bytes()
	return bodyBytes, nil
}

func wrapGzipRequest(r *resty.Request, gzippedBody []byte) {
	r.SetBody(gzippedBody)
	r.Header.Add("Content-Encoding", "gzip")
	r.Header.Add("Accept-Encoding", "gzip")
}

func sendMetricsBatch(metrics []*models.Metrics, host string, gzipRequest bool, keySHA string, ec chan error, semaphore chan struct{}) {
	logging.Log.Infoln("Sending metric batch")
	defer func() { <-semaphore }()
	client := resty.New()
	logging.Log.Infoln(host)
	url := "http://{host}/updates/"
	mJSON, err := json.Marshal(metrics)
	if err != nil {
		ec <- fmt.Errorf("error marshalling metric: %w", err)
		return
	}
	request := client.R().SetPathParam("host", host).
		SetHeader("Content-Type", "application/json")
	var body []byte
	if gzipRequest {
		body, err = wrapGzipBody(mJSON)
		if err != nil {
			ec <- fmt.Errorf("error gzipping metric: %w", err)
			return
		}
		wrapGzipRequest(request, body)
	} else {
		body = mJSON
		request.SetBody(mJSON)
	}
	if keySHA != "" {
		h := hmac.New(sha256.New, []byte(keySHA))
		hash := h.Sum(body)
		hashStr := base64.StdEncoding.EncodeToString(hash)
		request.SetHeader("Hash", hashStr)
	}
	res, err := trySendMetricsRetry(request, url)
	if err != nil {
		if res != nil {
			ec <- fmt.Errorf("error sending metric, response is not nil: %w, %d %s", err, res.StatusCode(), res.String())
			return
		} else {
			ec <- fmt.Errorf("error sending metric, response is nil: %w", err)
			return
		}
	}
	if res.StatusCode() != http.StatusOK {
		ec <- fmt.Errorf("status code not 200: %d %s", res.StatusCode(), res.String())
		return
	}
}

func trySendMetricsRetry(r *resty.Request, url string) (*resty.Response, error) {
	for timeout := range []int{1, 3, 5, -1} {
		res, err := r.Post(url)
		if err == nil {
			return res, err
		}
		if timeout < 0 {
			return res, err
		}
		if res != nil && res.StatusCode() == 0 {
			time.Sleep(time.Second * time.Duration(timeout))
		}
	}
	return nil, fmt.Errorf("timeout should end with <0")
}

func RunAgent(agentFlags *config.AgentFlags, runtimeMetrics []string, sendOnce bool, ignoreSendError bool) error {
	if err := logging.InitializeZapLogger("Info"); err != nil {
		return fmt.Errorf("error initializing logger: %w", err)
	}
	logging.Log.Infoln("agent started")
	collector := metrics.NewMetricCollector(runtimeMetrics)
	logging.Log.Infoln("Agent started with flags: ", *agentFlags)
	logging.Log.Infoln("Server Address:", agentFlags.Address)

	pollInterval := time.Second * time.Duration(agentFlags.PollInterval)

	log.Println("Starting metrics collection")
	errChannels := make([]chan error, 3)
	for i := 0; i < 3; i++ {
		errChannels[i] = make(chan error)
	}
	go CollectMetrics(collector, pollInterval, errChannels[0])
	go CollectAdditionalMetrics(collector, pollInterval, errChannels[1])
	go SendMetrics(collector, *agentFlags, sendOnce, ignoreSendError, errChannels[2])
	select {
	case err := <-errChannels[0]:
		return fmt.Errorf("metric collection error: %w", err)
	case err := <-errChannels[1]:
		return fmt.Errorf("additional metric collection error: %w", err)
	case err := <-errChannels[2]:
		if err != nil {
			return fmt.Errorf("metric send error: %w", err)
		} else {
			return nil
		}
	}
}
