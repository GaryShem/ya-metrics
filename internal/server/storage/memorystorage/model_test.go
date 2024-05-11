package memorystorage

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"
)

type MemStorageSuite struct {
	suite.Suite
	ms *MemStorage
}

func (s *MemStorageSuite) SetupTest() {
	s.ms = NewMemStorage()
}

func TestMemStorageSuite(t *testing.T) {
	suite.Run(t, new(MemStorageSuite))
}

func (s *MemStorageSuite) TestNewMemStorage() {
	ms := NewMemStorage()
	s.Require().NotNil(ms)
	s.True(reflect.DeepEqual(ms, s.ms))
}

func (s *MemStorageSuite) TestCounters() {
	value := int64(42)
	err := s.ms.UpdateCounter("foo", value)
	s.Require().NoError(err)
	counters, err := s.ms.GetCounters()
	s.Require().NoError(err)
	c, ok := counters["foo"]
	s.Require().True(ok)
	s.Equal(value, c.Value)
	foo, err := s.ms.GetCounter("foo")
	s.Require().NoError(err)
	s.Equal(value, foo.Value)

	err = s.ms.UpdateCounter("foo", value)
	s.Require().NoError(err)
	foo, err = s.ms.GetCounter("foo")
	s.Require().NoError(err)
	s.Equal(value*2, foo.Value)
}

func (s *MemStorageSuite) TestGauges() {
	value := float64(42)
	err := s.ms.UpdateGauge("foo", value)
	s.Require().NoError(err)

	gauges, err := s.ms.GetGauges()
	s.Require().NoError(err)
	c, ok := gauges["foo"]
	s.Require().True(ok)
	s.Equal(value, c.Value)

	foo, err := s.ms.GetGauge("foo")
	s.Require().NoError(err)
	s.Equal(value, foo.Value)
}
