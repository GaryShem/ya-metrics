package agent

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/GaryShem/ya-metrics.git/internal/server/handlers"
	"github.com/GaryShem/ya-metrics.git/internal/shared/logging"
	"github.com/GaryShem/ya-metrics.git/internal/shared/storage/memstorage"
	"github.com/GaryShem/ya-metrics.git/internal/shared/storage/models"
)

type AgentSuite struct {
	suite.Suite
	repo   models.Repository
	server *httptest.Server
	af     *AgentFlags
}

func (s *AgentSuite) SetupSuite() {
	s.repo = memstorage.NewMemStorage()
	router, err := handlers.MetricsRouter(s.repo)
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

	for _, m := range SupportedRuntimeMetrics() {
		g, err := s.repo.GetGauge(m)
		s.Require().NoError(err)
		s.Require().NotNil(g)
		s.Require().NotEqual(0, g.Value)
	}
	pc, err := s.repo.GetCounter("PollCount")
	s.Require().NoError(err)
	s.Require().NotNil(pc)
	s.Require().NotEqual(0, pc.Value)
}
