package shardedmap

import (
	"encoding/json"
)

// DefaultShardCount allows overwriting package default for New().
var (
	DefaultShardCount        uint              = 8             //nolint:gochecknoglobals
	DefaultShardProviderFunc ShardProviderFunc = NewMutexShard //nolint:gochecknoglobals
	DefaultKeyHashFunc       KeyHashFunc       = HashFnv1a64   //nolint:gochecknoglobals
)

// Map represents the sharded map.
type Map struct {
	shards            []Shard
	shardCount        uint
	shardProviderFunc ShardProviderFunc
	keyHashFunc       KeyHashFunc
}

// New creates a new sharded map.
func New(opts ...MapOption) *Map {
	m := &Map{} //nolint:exhaustivestruct
	m.applyDefaults()

	for _, opt := range opts {
		opt(m)
	}

	m.initShards()

	return m
}

func (m *Map) applyDefaults() {
	m.shardCount = DefaultShardCount
	m.shardProviderFunc = DefaultShardProviderFunc
	m.keyHashFunc = DefaultKeyHashFunc
}

func (m *Map) initShards() {
	m.shards = make([]Shard, m.shardCount)

	for j := 0; j < int(m.shardCount); j++ {
		m.shards[j] = m.shardProviderFunc()
	}
}

func (m *Map) getKeyHash(key string) uint {
	return m.keyHashFunc(key)
}

func (m *Map) calculateShardIndex(keyHash uint) uint {
	return keyHash % m.shardCount
}

func (m *Map) getKeyHashAndShardFromKey(key string) (keyHash uint, shard Shard) {
	keyHash = m.getKeyHash(key)
	shard = m.shards[m.calculateShardIndex(keyHash)]

	return
}

func (m *Map) forEachShard(cb func(shard Shard)) {
	for j := 0; j < int(m.shardCount); j++ {
		cb(m.shards[j])
	}
}

func (m *Map) RangeWithCallback(cb func(key string, value interface{}) interface{}) {
	m.forEachShard(func(shard Shard) {
		for keyHash, t := range shard.All() {
			newVal := cb(t.GetKey(), t.GetValue())
			if newVal != nil {
				shard.Set(keyHash, NewTuple(t.GetKey(), newVal))
			}
		}
	})
}

// TODO: Implement range with channel for easily use range
// func (m *Map) Range() chan i.ShardTuple {
// 	panic("not implemented")
// }

func (m *Map) Count() int {
	var count uint

	m.forEachShard(func(shard Shard) {
		count += shard.Count()
	})

	return int(count)
}

func (m *Map) All() map[string]interface{} {
	allKV := make(map[string]interface{})

	m.forEachShard(func(shard Shard) {
		for _, t := range shard.All() {
			allKV[t.GetKey()] = t.GetValue()
		}
	})

	return allKV
}

func (m *Map) Clear() {
	m.forEachShard(func(shard Shard) {
		shard.Clear()
	})
}

func (m *Map) Get(key string) interface{} {
	keyHash, shard := m.getKeyHashAndShardFromKey(key)
	val := shard.Get(keyHash)

	return val
}

func (m *Map) Set(key string, value interface{}) {
	keyHash, shard := m.getKeyHashAndShardFromKey(key)
	shard.Set(keyHash, NewTuple(key, value))
}

func (m *Map) Has(key string) bool {
	keyHash, shard := m.getKeyHashAndShardFromKey(key)
	has := shard.Has(keyHash)

	return has
}

func (m *Map) Remove(key string) {
	keyHash, shard := m.getKeyHashAndShardFromKey(key)
	shard.Remove(keyHash)
}

// UnmarshalJSON supports custom unmarshaling by implementing json.Unmarshaler interface.
func (m *Map) UnmarshalJSON(b []byte) error {
	if m.shards == nil && m.shardCount == 0 { // This is an empty Map
		m.applyDefaults()
		m.initShards()
	}

	flatMap := make(map[string]interface{})
	err := json.Unmarshal(b, &flatMap)

	if err != nil {
		return err
	}

	for k, v := range flatMap {
		m.Set(k, v)
	}

	return nil
}

// MarshalJSON supports custom marshaling by implementing json.Marshaler interface.
func (m *Map) MarshalJSON() ([]byte, error) {
	flatMap := m.All()

	return json.Marshal(flatMap)
}
