package shardedmap

// ShardProviderFunc defines a function that is used to create Shards while initializing a Map.
type ShardProviderFunc func() Shard

type ShardDataMap map[uint]ShardTuple

type ShardTuple interface {
	GetKey() string
	GetValue() interface{}
}

// Shard defines the interface that can be passed to map.
type Shard interface {

	// All returns all contained data as KVMap.
	All() ShardDataMap

	// Get returns a value from Collection.
	Get(uint) (interface{}, error)

	// Set sets a value to Collection.
	Set(uint, ShardTuple)

	// Has checks the existence of a key/value in Collection.
	Has(uint) bool

	// Remove removes an element by key from Collection.
	Remove(uint)

	// Count returns the count of elements in Collection.
	Count() uint

	// Clear resets all data in Collection.
	Clear()
}
