package handlers

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/GaryShem/ya-metrics.git/internal/shared/storage"
	"github.com/GaryShem/ya-metrics.git/internal/shared/storage/memstorage"
)

type MetricHandlerSuite struct {
	suite.Suite
	server *httptest.Server
	repo   storage.Repository
}

func (s *MetricHandlerSuite) BeforeTest(suiteName, testName string) {
	s.repo = memstorage.NewMemStorage()
	router, err := MetricsRouter(s.repo)
	if err != nil {
		panic(err)
	}
	s.server = httptest.NewServer(router)
}

func TestMetricHandlerSuite(t *testing.T) {
	suite.Run(t, new(MetricHandlerSuite))
}
