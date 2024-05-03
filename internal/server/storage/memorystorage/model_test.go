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
	s.ms.UpdateCounter("foo", value)
	counters := s.ms.GetCounters()
	c, ok := counters["foo"]
	s.Require().True(ok)
	s.Equal(value, c.Value)
	foo, err := s.ms.GetCounter("foo")
	s.Require().NoError(err)
	s.Equal(value, foo.Value)

	s.ms.UpdateCounter("foo", value)
	foo, err = s.ms.GetCounter("foo")
	s.Require().NoError(err)
	s.Equal(value*2, foo.Value)
}

func (s *MemStorageSuite) TestGauges() {
	value := float64(42)
	s.ms.UpdateGauge("foo", value)

	gauges := s.ms.GetGauges()
	c, ok := gauges["foo"]
	s.Require().True(ok)
	s.Equal(value, c.Value)

	foo, err := s.ms.GetGauge("foo")
	s.Require().NoError(err)
	s.Equal(value, foo.Value)
}
