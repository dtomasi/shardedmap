package shardedmap

// MapOption defines an option that can be set in Map constructor.
type MapOption func(m *Map)

// WithShardCount specifies the number of shards.
func WithShardCount(count int) MapOption {
	return func(m *Map) {
		m.shardCount = uint(count)
	}
}

// WithCustomShardProvider specifies the shard provider function to use.
func WithCustomShardProvider(provider ShardProviderFunc) MapOption {
	return func(m *Map) {
		m.shardProviderFunc = provider
	}
}

// WithCustomKeyHashFunc specifies the function for calculating the shard index for a given key.
func WithCustomKeyHashFunc(f KeyHashFunc) MapOption {
	return func(m *Map) {
		m.keyHashFunc = f
	}
}
