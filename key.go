package shardedmap

// KeyHashFunc defines a function that is used create a hash from a given key.
type KeyHashFunc func(key string) uint
