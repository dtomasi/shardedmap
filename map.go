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

func (m *Map) RangeWithCallback(cb func(key string, value interface{}) interface{}) {
	for _, shard := range m.shards {
		for keyHash, t := range shard.All() {
			newVal := cb(t.GetKey(), t.GetValue())
			if newVal != nil {
				shard.Set(keyHash, NewTuple(t.GetKey(), newVal))
			}
		}
	}
}

// Range allows iterating over a buffered data set.
func (m *Map) Range() <-chan ShardTuple {
	// The result channel which we return
	resChan := make(chan ShardTuple, m.Count())
	defer close(resChan)

	// A channel to which every shard if its is done
	shardDoneChan := make(chan bool, m.shardCount)
	defer close(shardDoneChan)

	// loop over shards
	for _, s := range m.shards {
		// Fetch all data in a separate goroutine
		go func(shard Shard, doneChan chan bool) {
			// Push results to resChan
			for _, v := range shard.All() {
				resChan <- v
			}
			// Notify that we are done here
			doneChan <- true
		}(s, shardDoneChan)
	}

	var shardsDoneCount uint

	// Wait for all to finish
	for { // nolint:gosimple
		select {
		case <-shardDoneChan:
			shardsDoneCount++
			if shardsDoneCount == m.shardCount {
				return resChan
			}
		}
	}
}

// Count returns the count of all elements across all shards.
func (m *Map) Count() int {
	// The result channel which we return
	countChan := make(chan uint)
	defer close(countChan)

	// A channel to which every shard if its is done
	shardDoneChan := make(chan bool, m.shardCount)
	defer close(shardDoneChan)

	// loop over shards
	for _, s := range m.shards {
		// Fetch all data in a separate goroutine
		go func(shard Shard, doneChan chan bool) {
			// Push results to resChan
			countChan <- shard.Count()
			// Notify that we are done here
			doneChan <- true
		}(s, shardDoneChan)
	}

	var (
		totalCount      uint
		shardsDoneCount uint
	)

	// Wait for all to finish
	for {
		select {
		case countFromShard := <-countChan:
			totalCount += countFromShard

		case <-shardDoneChan:
			shardsDoneCount++
			if shardsDoneCount == m.shardCount {
				return int(totalCount)
			}
		}
	}
}

// All returns a flat map of all keys and values across all shards.
func (m *Map) All() map[string]interface{} {
	allData := make(map[string]interface{})

	for valTuple := range m.Range() {
		allData[valTuple.GetKey()] = valTuple.GetValue()
	}

	return allData
}

// Clear clears all data across all shards.
func (m *Map) Clear() {
	for _, shard := range m.shards {
		shard.Clear()
	}
}

// Get returns the value for given key or an error.
func (m *Map) Get(key string) (interface{}, error) {
	keyHash, shard := m.getKeyHashAndShardFromKey(key)
	val, err := shard.Get(keyHash)

	if err != nil {
		return nil, err
	}

	return val, nil
}

// MustGet returns the value for a given key or nil.
func (m *Map) MustGet(key string) interface{} {
	val, _ := m.Get(key)

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
