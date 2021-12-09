package shardedmap

import (
	"errors"
	"sync"
	"sync/atomic"
)

// NewAtomicShard creates a new AtomicShard.
func NewAtomicShard() Shard {
	s := new(AtomicShard)
	s.data.Store(make(ShardDataMap))

	return s
}

// AtomicShard represents a shard used in Map.
type AtomicShard struct {
	mu   sync.RWMutex
	data atomic.Value
}

func (s *AtomicShard) getValueMap() ShardDataMap {
	m1 := s.data.Load().(ShardDataMap) // nolint:forcetypeassert

	return m1
}

// All see: interfaces.Shard.
func (s *AtomicShard) All() ShardDataMap {
	return s.getValueMap()
}

// Get see: interfaces.Shard.
func (s *AtomicShard) Get(key uint) (interface{}, error) {
	tuple, ok := s.getValueMap()[key]
	if !ok {
		return nil, errors.New("not found") // nolint
	}

	return tuple.GetValue(), nil
}

// Set see: interfaces.Shard.
func (s *AtomicShard) Set(key uint, value ShardTuple) {
	m1 := s.getValueMap()
	s.mu.Lock()
	defer s.mu.Unlock()

	m1[key] = value

	s.data.Store(m1)
}

// Has see: interfaces.Shard.
func (s *AtomicShard) Has(key uint) bool {
	if _, ok := s.getValueMap()[key]; ok {
		return true
	}

	return false
}

// Remove see: interfaces.Shard.
func (s *AtomicShard) Remove(key uint) {
	m1 := s.getValueMap()
	delete(m1, key)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data.Store(m1)
}

// Count see: interfaces.Shard.
func (s *AtomicShard) Count() uint {
	return uint(len(s.getValueMap()))
}

// Clear see: interfaces.Shard.
func (s *AtomicShard) Clear() {
	s.data.Store(make(ShardDataMap))
}
