module GoLevelBadgerRocksPebbleBenchmark

go 1.16

require (
	github.com/cockroachdb/pebble v0.0.0-20210716215258-b8b89851512d
	github.com/dgraph-io/badger/v3 v3.2103.1
	github.com/facebookgo/ensure v0.0.0-20200202191622-63f1cf65ac4c // indirect
	github.com/facebookgo/stack v0.0.0-20160209184415-751773369052 // indirect
	github.com/facebookgo/subset v0.0.0-20200203212716-c811ad88dec4 // indirect
	github.com/syndtr/goleveldb v1.0.0
	github.com/tecbot/gorocksdb v0.0.0-20191217155057-f0fad39f321c
)

replace github.com/dgraph-io/badger/v3 => ../../../github.com/dgraph-io/badger/
