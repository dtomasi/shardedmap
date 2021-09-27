package shardedmap_test

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/dtomasi/shardedmap"
	"github.com/stretchr/testify/suite"
	"math/rand"
	"reflect"
)

func pickRandomKeyFromDataSet(m map[string]interface{}) string {
	keys := reflect.ValueOf(m).MapKeys()

	return keys[rand.Intn(len(keys))].String() //nolint:gosec
}

type ShardTestSuite struct {
	suite.Suite
	testDataSet map[string]interface{}
	instance    shardedmap.Shard
}

func NewShardTestSuite(shardInstance shardedmap.Shard) *ShardTestSuite {
	return &ShardTestSuite{instance: shardInstance} //nolint:exhaustivestruct
}

func (s *ShardTestSuite) SetupTest() {
	// Clear map to get a new instance here
	s.instance.Clear()

	// Generate test data
	s.testDataSet = gofakeit.Map()
	for k, v := range s.testDataSet {
		s.instance.Set(shardedmap.HashFnv1a64(k), shardedmap.NewTuple(k, v))
	}
}

func (s *ShardTestSuite) TestGet() {
	k := pickRandomKeyFromDataSet(s.testDataSet)
	v, err := s.instance.Get(shardedmap.HashFnv1a64(k))
	s.NoError(err)
	s.Equal(s.testDataSet[k], v)
}

func (s *ShardTestSuite) TestHas() {
	k := pickRandomKeyFromDataSet(s.testDataSet)
	s.True(s.instance.Has(shardedmap.HashFnv1a64(k)))
}

func (s *ShardTestSuite) TestSet() {
	keyHash := shardedmap.HashFnv1a64("key")
	s.instance.Set(keyHash, shardedmap.NewTuple("key", "value"))
	s.True(s.instance.Has(keyHash))
}

func (s *ShardTestSuite) TestCount() {
	s.Equal(len(s.testDataSet), int(s.instance.Count()))
}

func (s *ShardTestSuite) TestRemove() {
	keyHash := shardedmap.HashFnv1a64(pickRandomKeyFromDataSet(s.testDataSet))
	s.instance.Remove(keyHash)
	s.False(s.instance.Has(keyHash))
}

func (s *ShardTestSuite) TestClear() {
	s.instance.Clear()
	s.Equal(0, int(s.instance.Count()))
}

func (s *ShardTestSuite) TestAll() {
	// Convert back into map[string]interface{}
	shardData := make(map[string]interface{}, len(s.testDataSet))
	for _, v := range s.instance.All() {
		shardData[v.GetKey()] = v.GetValue()
	}

	s.Equal(s.testDataSet, shardData)
	s.Equal(len(s.testDataSet), len(shardData))
}
