package agent

import (
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/GaryShem/ya-metrics.git/internal/server"
	"github.com/GaryShem/ya-metrics.git/internal/shared/storage"
)

func TestRunAgent(t *testing.T) {
	reportInteral := 2
	pollInterval := 1
	router, err := server.MetricsRouter(storage.NewMemStorage())
	if err != nil {
		t.Fatal(err)
	}
	ts := httptest.NewServer(router)
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
