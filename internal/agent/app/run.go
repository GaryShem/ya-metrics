package app

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/GaryShem/ya-metrics.git/internal/agent/config"
	"github.com/GaryShem/ya-metrics.git/internal/agent/encryption"
	"github.com/GaryShem/ya-metrics.git/internal/agent/metrics"
	"github.com/GaryShem/ya-metrics.git/internal/agent/metrics/sender"
	"github.com/GaryShem/ya-metrics.git/internal/shared/logging"
)

func RunAgent(agentFlags *config.AgentFlags, runtimeMetrics []string, sendOnce bool) error {
	if err := logging.InitializeZapLogger("Info"); err != nil {
		return fmt.Errorf("error initializing logger: %w", err)
	}
	if agentFlags.CryptoKey != "" {
		err := encryption.InitEncryptor(agentFlags.CryptoKey)
		if err != nil {
			return fmt.Errorf("error initializing encryptor: %w", err)
		}
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
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()
	go sender.CollectMetrics(ctx, collector, pollInterval, errChannels[0])
	go sender.CollectAdditionalMetrics(ctx, collector, pollInterval, errChannels[1])
	go sender.SendMetrics(ctx, collector, *agentFlags, sendOnce, errChannels[2])
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
