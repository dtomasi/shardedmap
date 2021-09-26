# shardedmap

[![Go Reference](https://pkg.go.dev/badge/github.com/dtomasi/shardedmap.svg)](https://pkg.go.dev/github.com/dtomasi/shardedmap)
[![CodeFactor](https://www.codefactor.io/repository/github/dtomasi/shardedmap/badge)](https://www.codefactor.io/repository/github/dtomasi/shardedmap)
![Go Unit Tests](https://github.com/dtomasi/shardedmap/actions/workflows/build.yml/badge.svg)
![Go Benchmarks](https://github.com/dtomasi/shardedmap/actions/workflows/benchmark_cob.yml/badge.svg)
![CodeQL](https://github.com/dtomasi/shardedmap/actions/workflows/codeql-analysis.yml/badge.svg)
[![codecov](https://codecov.io/gh/dtomasi/shardedmap/branch/main/graph/badge.svg?token=9fqDqF2uxF)](https://codecov.io/gh/dtomasi/shardedmap)

This project contains a threadsafe sharded map implementation in Go.

## Installation

    go get github.com/dtomasi/shardedmap

## Usage

TBD

## Notice

The idea of splitting maps into shards to solve parallel access issues is not new. I borrowed some ideas from here:

https://github.com/orcaman/concurrent-map

https://github.com/allegro/bigcache

## LICENCE

see [LICENCE](LICENSE)
