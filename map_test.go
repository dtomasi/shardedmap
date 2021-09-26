package shardedmap_test

import (
	"github.com/dtomasi/shardedmap"
	"github.com/stretchr/testify/suite"
	"testing"
)

var mapTestMatrix = []struct { //nolint:gochecknoglobals
	shardCount    int
	shardProvider shardedmap.ShardProviderFunc
	keyHashFunc   shardedmap.KeyHashFunc
}{
	// Shard Count 1
	{1, shardedmap.NewMutexShard, shardedmap.HashFnv1a64},
	{1, shardedmap.NewMutexShard, shardedmap.HashFnv1a32},
	{1, shardedmap.NewAtomicShard, shardedmap.HashFnv1a64},
	{1, shardedmap.NewAtomicShard, shardedmap.HashFnv1a32},

	// Shard Count 8
	{8, shardedmap.NewMutexShard, shardedmap.HashFnv1a64},
	{8, shardedmap.NewMutexShard, shardedmap.HashFnv1a32},
	{8, shardedmap.NewAtomicShard, shardedmap.HashFnv1a64},
	{8, shardedmap.NewAtomicShard, shardedmap.HashFnv1a32},

	// Shard Count 32
	{32, shardedmap.NewMutexShard, shardedmap.HashFnv1a64},
	{32, shardedmap.NewMutexShard, shardedmap.HashFnv1a32},
	{32, shardedmap.NewAtomicShard, shardedmap.HashFnv1a64},
	{32, shardedmap.NewAtomicShard, shardedmap.HashFnv1a32},

	// Shard Count 128
	{128, shardedmap.NewMutexShard, shardedmap.HashFnv1a64},
	{128, shardedmap.NewMutexShard, shardedmap.HashFnv1a32},
	{128, shardedmap.NewAtomicShard, shardedmap.HashFnv1a64},
	{128, shardedmap.NewAtomicShard, shardedmap.HashFnv1a32},
}

func TestMapRunSuiteMatrix(t *testing.T) {
	for _, testSetup := range mapTestMatrix {
		suite.Run(t, NewMapTestSuite(
			shardedmap.New(
				shardedmap.WithShardCount(testSetup.shardCount),
				shardedmap.WithCustomShardProvider(testSetup.shardProvider),
				shardedmap.WithCustomKeyHashFunc(testSetup.keyHashFunc),
			),
		))
	}
}
