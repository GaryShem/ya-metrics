package agent

import (
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/GaryShem/ya-metrics.git/internal/server"
	"github.com/GaryShem/ya-metrics.git/internal/shared/storage"
	"github.com/GaryShem/ya-metrics.git/internal/shared/storage/metrics"
)

func TestRunAgent(t *testing.T) {
	reportInteral := 2
	pollInterval := 1

	ts := httptest.NewServer(server.MetricsRouter(
		&storage.MemStorage{
			GaugeMetrics:   map[string]*metrics.Gauge{},
			CounterMetrics: map[string]*metrics.Counter{},
		},
	))
	defer ts.Close()
	serverURLSlice := strings.Split(ts.URL, ":")
	serverPort := serverURLSlice[len(serverURLSlice)-1]
	serverAddress := fmt.Sprintf("127.0.0.1:%s", serverPort)
	type args struct {
		af       *AgentFlags
		sendOnce bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "Send Stats Once",
			args: args{
				af: &AgentFlags{
					Address:        &serverAddress,
					ReportInterval: &reportInteral,
					PollInterval:   &pollInterval,
				},
				sendOnce: true,
			},
			wantErr: require.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t,
				RunAgent(tt.args.af, tt.args.sendOnce),
				fmt.Sprintf("RunAgent(%v, %v)", tt.args.af, tt.args.sendOnce))
		})
	}
}
