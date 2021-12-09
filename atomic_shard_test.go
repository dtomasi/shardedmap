package shardedmap_test

import (
	"github.com/dtomasi/shardedmap"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestAtomicShardTestsInSuite(t *testing.T) {
	suite.Run(t, NewShardTestSuite(shardedmap.NewAtomicShard()))
}
