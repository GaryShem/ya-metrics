package agent

import (
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/GaryShem/ya-metrics.git/internal/server/handlers"
	"github.com/GaryShem/ya-metrics.git/internal/shared/logging"
	"github.com/GaryShem/ya-metrics.git/internal/shared/storage"
	"github.com/GaryShem/ya-metrics.git/internal/shared/storage/memstorage"
)

type AgentSuite struct {
	suite.Suite
	repo   storage.Repository
	server *httptest.Server
	af     *AgentFlags
}

func (s *AgentSuite) SetupSuite() {
	s.repo = memstorage.NewMemStorage()
	router, err := handlers.MetricsRouter(memstorage.NewMemStorage())
	if err != nil {
		panic(err)
	}
	s.server = httptest.NewServer(router)
	logging.Log.Infoln("server url:", s.server.URL)
	serverURLSlice := strings.Split(s.server.URL, ":")
	//serverIP := serverURLSlice[len(serverURLSlice)-2]
	serverPort := serverURLSlice[len(serverURLSlice)-1]
	logging.Log.Infoln(s.server.URL, serverPort)
	serverAddress, _ := strings.CutPrefix(s.server.URL, "http://")
	logging.Log.Infoln("server address", serverAddress)
	reportInterval := 2
	pollInterval := 1
	s.af = &AgentFlags{
		Address:        &serverAddress,
		ReportInterval: &reportInterval,
		PollInterval:   &pollInterval,
	}
}

func (s *AgentSuite) TearDownSuite() {
	s.server.Close()
}

func TestAgentSuite(t *testing.T) {
	suite.Run(t, new(AgentSuite))
}

func (s *AgentSuite) TestAgentMetrics() {
	err := RunAgent(s.af, true)
	s.Require().NoError(err)

	//query := s.server.URL + "/"
	//client :=
}

func TestRunAgent(t *testing.T) {
	reportInteral := 2
	pollInterval := 1
	router, err := handlers.MetricsRouter(memstorage.NewMemStorage())
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
