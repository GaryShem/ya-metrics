package server

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/GaryShem/ya-metrics.git/internal/shared/storage"
)

type MetricHandlerSuite struct {
	suite.Suite
	server *httptest.Server
	repo   storage.Repository
}

func (s *MetricHandlerSuite) SetupSuite() {
	s.repo = storage.NewMemStorage()
	router, err := MetricsRouter(s.repo)
	if err != nil {
		panic(err)
	}
	s.server = httptest.NewServer(router)
}

func TestMetricHandlerSuite(t *testing.T) {
	suite.Run(t, new(MetricHandlerSuite))
}
