package shardedmap

import (
	"sync"
)

// NewMutexShard creates a new Shard.
func NewMutexShard() Shard {
	return &MutexShard{ //nolint:exhaustivestruct
		data: make(ShardDataMap),
	}
}

// MutexShard represents a shard used in Map.
type MutexShard struct {
	mu   sync.RWMutex
	data ShardDataMap
}

// All see: interfaces.Collection.
func (s *MutexShard) All() ShardDataMap {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.data
}

// Get see: interfaces.Collection.
func (s *MutexShard) Get(key uint) interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.data[key].GetValue()
}

// Set see: interfaces.Collection.
func (s *MutexShard) Set(key uint, value ShardTuple) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

// Has see: interfaces.Collection.
func (s *MutexShard) Has(key uint) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, ok := s.data[key]; ok {
		return true
	}

	return false
}

// Remove see: interfaces.Collection.
func (s *MutexShard) Remove(key uint) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
}

// Count see: interfaces.Collection.
func (s *MutexShard) Count() uint {
	// TODO: Check if builtin len function requires a read lock here
	s.mu.RLock()
	defer s.mu.RUnlock()

	return uint(len(s.data))
}

// Clear see: interfaces.Collection.
func (s *MutexShard) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = make(ShardDataMap)
}
