package shardedmap_test

import (
	"encoding/json"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/dtomasi/shardedmap"
	"github.com/stretchr/testify/suite"
	"reflect"
)

type MapTestSuite struct {
	suite.Suite
	testDataSet map[string]interface{}
	instance    *shardedmap.Map
}

func NewMapTestSuite(mapInstance *shardedmap.Map) *MapTestSuite {
	return &MapTestSuite{instance: mapInstance} //nolint:exhaustivestruct
}

func (s *MapTestSuite) SetupTest() { // Generate test data
	// Clear map to get a new instance here
	s.instance.Clear()

	// Generate test data
	s.testDataSet = gofakeit.Map()
	for k, v := range s.testDataSet {
		s.instance.Set(k, v)
	}
}

func (s *MapTestSuite) TestGet() {
	k := pickRandomKeyFromDataSet(s.testDataSet)
	v := s.instance.Get(k)
	s.Equal(s.testDataSet[k], v)
}

func (s *MapTestSuite) TestHas() {
	k := pickRandomKeyFromDataSet(s.testDataSet)
	s.True(s.instance.Has(k))
}

func (s *MapTestSuite) TestSet() {
	s.instance.Set("key", "value")
	s.True(s.instance.Has("key"))
}

func (s *MapTestSuite) TestCount() {
	s.Equal(len(s.testDataSet), s.instance.Count())
}

func (s *MapTestSuite) TestRemove() {
	key := pickRandomKeyFromDataSet(s.testDataSet)
	s.instance.Remove(key)
	s.False(s.instance.Has(key))
}

func (s *MapTestSuite) TestClear() {
	s.instance.Clear()
	s.Equal(0, s.instance.Count())
}

func (s *MapTestSuite) TestAll() {
	// Convert back into map[string]interface{}
	shardData := make(map[string]interface{}, len(s.testDataSet))
	for k, v := range s.instance.All() {
		shardData[k] = v

		// Best effort here ... Value could be a map or something that we cannot compare
		if reflect.TypeOf(v).Comparable() {
			s.Equal(s.testDataSet[k], v)
		}
	}

	s.Equal(s.testDataSet, shardData)
	s.Equal(len(s.testDataSet), len(shardData))
}

func (s *MapTestSuite) TestRangeWithCallback() {
	s.instance.RangeWithCallback(func(key string, value interface{}) interface{} {
		// Best effort here ... Value could be a map or something that we cannot compare
		if reflect.TypeOf(value).Comparable() {
			s.Equal(s.testDataSet[key], value)
		}

		return "empty"
	})

	// As we have overwritten each value with "empty" ... Let´s check a random key
	randomKey := pickRandomKeyFromDataSet(s.testDataSet)
	s.Equal("empty", s.instance.Get(randomKey))
}

func (s *MapTestSuite) TestJsonMarshalAndUnmarshal() {
	// MarshalJSON
	jsonBytes, marshalErr := json.Marshal(s.instance)
	s.NoError(marshalErr)
	s.True(json.Valid(jsonBytes))

	// UnmarshalJSON
	var m *shardedmap.Map // Note: This will create a map with defaults set
	unmarshalErr := json.Unmarshal(jsonBytes, &m)
	s.NoError(unmarshalErr)
	s.Equal(len(s.testDataSet), m.Count())
}
