package http

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/GaryShem/ya-metrics.git/internal/server/storage/memorystorage"
	"github.com/GaryShem/ya-metrics.git/internal/server/storage/repository"
)

type MetricHandlerSuite struct {
	suite.Suite
	server *httptest.Server
	repo   repository.Repository
}

func (s *MetricHandlerSuite) BeforeTest(suiteName, testName string) {
	s.repo = memorystorage.NewMemStorage()
	router, err := MetricsRouter(s.repo, false)
	if err != nil {
		panic(err)
	}
	s.server = httptest.NewServer(router)
}

func TestMetricHandlerSuite(t *testing.T) {
	suite.Run(t, new(MetricHandlerSuite))
}
