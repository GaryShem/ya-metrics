package agent

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/GaryShem/ya-metrics.git/internal/server/handlers"
	"github.com/GaryShem/ya-metrics.git/internal/server/storage/memorystorage"
	"github.com/GaryShem/ya-metrics.git/internal/shared/logging"
	"github.com/GaryShem/ya-metrics.git/internal/shared/storage/models"
)

type AgentSuite struct {
	suite.Suite
	repo   models.Repository
	server *httptest.Server
	af     *AgentFlags
}

func (s *AgentSuite) SetupSuite() {
	s.repo = memorystorage.NewMemStorage()
	router, err := handlers.MetricsRouter(s.repo)
	if err != nil {
		panic(err)
	}
	s.server = httptest.NewServer(router)
	logging.Log.Infoln("server url:", s.server.URL)
	serverURLSlice := strings.Split(s.server.URL, ":")
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
	err := RunAgent(s.af, SupportedRuntimeMetrics(),
		true, false, false)
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

func (s *AgentSuite) TestAgentGzip() {
	metrics := []string{"Alloc"}
	err := RunAgent(s.af, metrics,
		true, false, true)
	s.Require().NoError(err)

	for _, m := range metrics {
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
