package shardedmap_test

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/dtomasi/shardedmap"
	"math/rand"
	"runtime"
	"testing"
	"time"
)

const sleepAfterBenchmarkDuration = time.Second * 1

//nolint:gochecknoinits
func init() {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())
}

func runSequentialBenchmarkSet(
	b *testing.B,
	shardCount int,
	shardProvider shardedmap.ShardProviderFunc,
	keyHashFunc shardedmap.KeyHashFunc,
) {
	b.Helper()

	instance := shardedmap.New(
		shardedmap.WithShardCount(shardCount),
		shardedmap.WithCustomShardProvider(shardProvider),
		shardedmap.WithCustomKeyHashFunc(keyHashFunc),
	)
	testData := gofakeit.Map()

	for k, v := range testData {
		instance.Set(k, v)
	}

	randomKey := pickRandomKeyFromDataSet(testData)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		instance.Set(randomKey, testData[randomKey])
	}

	// Give go some time to breath
	b.StopTimer()
	runtime.GC()
	time.Sleep(sleepAfterBenchmarkDuration)
}

func Benchmark_ShardedMap_Sequential_ShardCount_1__Provider_Mutex__Hash_32__Set(b *testing.B) {
	runSequentialBenchmarkSet(b, 1, shardedmap.NewMutexShard, shardedmap.HashFnv1a32)
}

func Benchmark_ShardedMap_Sequential_ShardCount_1__Provider_Mutex__Hash_64__Set(b *testing.B) {
	runSequentialBenchmarkSet(b, 1, shardedmap.NewMutexShard, shardedmap.HashFnv1a64)
}

func Benchmark_ShardedMap_Sequential_ShardCount_1__Provider_Atomic__Hash_32__Set(b *testing.B) {
	runSequentialBenchmarkSet(b, 1, shardedmap.NewAtomicShard, shardedmap.HashFnv1a32)
}

func Benchmark_ShardedMap_Sequential_ShardCount_1__Provider_Atomic__Hash_64__Set(b *testing.B) {
	runSequentialBenchmarkSet(b, 1, shardedmap.NewAtomicShard, shardedmap.HashFnv1a64)
}

func Benchmark_ShardedMap_Sequential_ShardCount_32__Provider_Mutex__Hash_32__Set(b *testing.B) {
	runSequentialBenchmarkSet(b, 32, shardedmap.NewMutexShard, shardedmap.HashFnv1a32)
}

func Benchmark_ShardedMap_Sequential_ShardCount_32__Provider_Mutex__Hash_64__Set(b *testing.B) {
	runSequentialBenchmarkSet(b, 32, shardedmap.NewMutexShard, shardedmap.HashFnv1a64)
}

func Benchmark_ShardedMap_Sequential_ShardCount_32__Provider_Atomic__Hash_32__Set(b *testing.B) {
	runSequentialBenchmarkSet(b, 32, shardedmap.NewAtomicShard, shardedmap.HashFnv1a32)
}

func Benchmark_ShardedMap_Sequential_ShardCount_32__Provider_Atomic__Hash_64__Set(b *testing.B) {
	runSequentialBenchmarkSet(b, 32, shardedmap.NewAtomicShard, shardedmap.HashFnv1a64)
}

func runParallelBenchmarkSet(
	b *testing.B,
	shardCount int,
	shardProvider shardedmap.ShardProviderFunc,
	keyHashFunc shardedmap.KeyHashFunc,
) {
	b.Helper()

	instance := shardedmap.New(
		shardedmap.WithShardCount(shardCount),
		shardedmap.WithCustomShardProvider(shardProvider),
		shardedmap.WithCustomKeyHashFunc(keyHashFunc),
	)
	testData := gofakeit.Map()

	for k, v := range testData {
		instance.Set(k, v)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		randomKey := pickRandomKeyFromDataSet(testData)
		for pb.Next() {
			instance.Set(randomKey, testData[randomKey])
		}
	})
	// Give go some time to breath
	b.StopTimer()
	runtime.GC()
	time.Sleep(sleepAfterBenchmarkDuration)
}

func Benchmark_ShardedMap_Parallel_ShardCount_1__Provider_Mutex__Hash_32__Set(b *testing.B) {
	runParallelBenchmarkSet(b, 1, shardedmap.NewMutexShard, shardedmap.HashFnv1a32)
}

func Benchmark_ShardedMap_Parallel_ShardCount_1__Provider_Mutex__Hash_64__Set(b *testing.B) {
	runParallelBenchmarkSet(b, 1, shardedmap.NewMutexShard, shardedmap.HashFnv1a64)
}

func Benchmark_ShardedMap_Parallel_ShardCount_1__Provider_Atomic__Hash_32__Set(b *testing.B) {
	runParallelBenchmarkSet(b, 1, shardedmap.NewAtomicShard, shardedmap.HashFnv1a32)
}

func Benchmark_ShardedMap_Parallel_ShardCount_1__Provider_Atomic__Hash_64__Set(b *testing.B) {
	runParallelBenchmarkSet(b, 1, shardedmap.NewAtomicShard, shardedmap.HashFnv1a64)
}

func Benchmark_ShardedMap_Parallel_ShardCount_32__Provider_Mutex__Hash_32__Set(b *testing.B) {
	runParallelBenchmarkSet(b, 32, shardedmap.NewMutexShard, shardedmap.HashFnv1a32)
}

func Benchmark_ShardedMap_Parallel_ShardCount_32__Provider_Mutex__Hash_64__Set(b *testing.B) {
	runParallelBenchmarkSet(b, 32, shardedmap.NewMutexShard, shardedmap.HashFnv1a64)
}

func Benchmark_ShardedMap_Parallel_ShardCount_32__Provider_Atomic__Hash_32__Set(b *testing.B) {
	runParallelBenchmarkSet(b, 32, shardedmap.NewAtomicShard, shardedmap.HashFnv1a32)
}

func Benchmark_ShardedMap_Parallel_ShardCount_32__Provider_Atomic__Hash_64__Set(b *testing.B) {
	runParallelBenchmarkSet(b, 32, shardedmap.NewAtomicShard, shardedmap.HashFnv1a64)
}

func runSequentialBenchmarkGet(
	b *testing.B,
	shardCount int,
	shardProvider shardedmap.ShardProviderFunc,
	keyHashFunc shardedmap.KeyHashFunc,
) {
	b.Helper()

	instance := shardedmap.New(
		shardedmap.WithShardCount(shardCount),
		shardedmap.WithCustomShardProvider(shardProvider),
		shardedmap.WithCustomKeyHashFunc(keyHashFunc),
	)
	testData := gofakeit.Map()

	for k, v := range testData {
		instance.Set(k, v)
	}

	randomKey := pickRandomKeyFromDataSet(testData)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		v := instance.Get(randomKey)
		if v == nil {
			b.FailNow()
		}
	}
	// Give go some time to breath
	b.StopTimer()
	runtime.GC()
	time.Sleep(sleepAfterBenchmarkDuration)
}

func Benchmark_ShardedMap_Sequential_ShardCount_1__Provider_Mutex__Hash_32__Get(b *testing.B) {
	runSequentialBenchmarkGet(b, 1, shardedmap.NewMutexShard, shardedmap.HashFnv1a32)
}

func Benchmark_ShardedMap_Sequential_ShardCount_1__Provider_Mutex__Hash_64__Get(b *testing.B) {
	runSequentialBenchmarkGet(b, 1, shardedmap.NewMutexShard, shardedmap.HashFnv1a64)
}

func Benchmark_ShardedMap_Sequential_ShardCount_1__Provider_Atomic__Hash_32__Get(b *testing.B) {
	runSequentialBenchmarkGet(b, 1, shardedmap.NewAtomicShard, shardedmap.HashFnv1a32)
}

func Benchmark_ShardedMap_Sequential_ShardCount_1__Provider_Atomic__Hash_64__Get(b *testing.B) {
	runSequentialBenchmarkGet(b, 1, shardedmap.NewAtomicShard, shardedmap.HashFnv1a64)
}

func Benchmark_ShardedMap_Sequential_ShardCount_32__Provider_Mutex__Hash_32__Get(b *testing.B) {
	runSequentialBenchmarkGet(b, 32, shardedmap.NewMutexShard, shardedmap.HashFnv1a32)
}

func Benchmark_ShardedMap_Sequential_ShardCount_32__Provider_Mutex__Hash_64__Get(b *testing.B) {
	runSequentialBenchmarkGet(b, 32, shardedmap.NewMutexShard, shardedmap.HashFnv1a64)
}

func Benchmark_ShardedMap_Sequential_ShardCount_32__Provider_Atomic__Hash_32__Get(b *testing.B) {
	runSequentialBenchmarkGet(b, 32, shardedmap.NewAtomicShard, shardedmap.HashFnv1a32)
}

func Benchmark_ShardedMap_Sequential_ShardCount_32__Provider_Atomic__Hash_64__Get(b *testing.B) {
	runSequentialBenchmarkGet(b, 32, shardedmap.NewAtomicShard, shardedmap.HashFnv1a64)
}

func runParallelBenchmarkGet(
	b *testing.B,
	shardCount int,
	shardProvider shardedmap.ShardProviderFunc,
	keyHashFunc shardedmap.KeyHashFunc,
) {
	b.Helper()

	instance := shardedmap.New(
		shardedmap.WithShardCount(shardCount),
		shardedmap.WithCustomShardProvider(shardProvider),
		shardedmap.WithCustomKeyHashFunc(keyHashFunc),
	)
	testData := gofakeit.Map()

	for k, v := range testData {
		instance.Set(k, v)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		randomKey := pickRandomKeyFromDataSet(testData)
		for pb.Next() {
			v := instance.Get(randomKey)
			if v == nil {
				b.FailNow()
			}
		}
	})
	// Give go some time to breath
	b.StopTimer()
	runtime.GC()
	time.Sleep(sleepAfterBenchmarkDuration)
}

func Benchmark_ShardedMap_Parallel_ShardCount_1__Provider_Mutex__Hash_32__Get(b *testing.B) {
	runParallelBenchmarkGet(b, 1, shardedmap.NewMutexShard, shardedmap.HashFnv1a32)
}

func Benchmark_ShardedMap_Parallel_ShardCount_1__Provider_Mutex__Hash_64__Get(b *testing.B) {
	runParallelBenchmarkGet(b, 1, shardedmap.NewMutexShard, shardedmap.HashFnv1a64)
}

func Benchmark_ShardedMap_Parallel_ShardCount_1__Provider_Atomic__Hash_32__Get(b *testing.B) {
	runParallelBenchmarkGet(b, 1, shardedmap.NewAtomicShard, shardedmap.HashFnv1a32)
}

func Benchmark_ShardedMap_Parallel_ShardCount_1__Provider_Atomic__Hash_64__Get(b *testing.B) {
	runParallelBenchmarkGet(b, 1, shardedmap.NewAtomicShard, shardedmap.HashFnv1a64)
}

func Benchmark_ShardedMap_Parallel_ShardCount_32__Provider_Mutex__Hash_32__Get(b *testing.B) {
	runParallelBenchmarkGet(b, 32, shardedmap.NewMutexShard, shardedmap.HashFnv1a32)
}

func Benchmark_ShardedMap_Parallel_ShardCount_32__Provider_Mutex__Hash_64__Get(b *testing.B) {
	runParallelBenchmarkGet(b, 32, shardedmap.NewMutexShard, shardedmap.HashFnv1a64)
}

func Benchmark_ShardedMap_Parallel_ShardCount_32__Provider_Atomic__Hash_32__Get(b *testing.B) {
	runParallelBenchmarkGet(b, 32, shardedmap.NewAtomicShard, shardedmap.HashFnv1a32)
}

func Benchmark_ShardedMap_Parallel_ShardCount_32__Provider_Atomic__Hash_64__Get(b *testing.B) {
	runParallelBenchmarkGet(b, 32, shardedmap.NewAtomicShard, shardedmap.HashFnv1a64)
}
