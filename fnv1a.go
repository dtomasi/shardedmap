package shardedmap

const (
	// 32-Bit.
	fnvOffset32 = 2166136261
	fnvPrime32  = 16777619

	// 64-Bit.
	fnvOffset64 = 14695981039346656037
	fnvPrime64  = 1099511628211
)

func HashFnv1a32(key string) uint {
	var _hash uint64 = fnvOffset32
	for j := 0; j < len(key); j++ {
		_hash ^= uint64(key[j])
		_hash *= fnvPrime32
	}

	return uint(_hash)
}

func HashFnv1a64(key string) uint {
	var _hash uint64 = fnvOffset64
	for j := 0; j < len(key); j++ {
		_hash ^= uint64(key[j])
		_hash *= fnvPrime64
	}

	return uint(_hash)
}
