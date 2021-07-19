ulimit -a
cd /home/gopath/GoLevelBadgerRocksPebbleBenchmark/BadgerDB
go run main.go
cd /tmp/badgertest-0/dbbench
du --total
cd /home/gopath/GoLevelBadgerRocksPebbleBenchmark
./clean.sh
cd /home/gopath/GoLevelBadgerRocksPebbleBenchmark/PebbleDB
go run main.go
cd /tmp/pebbledbtest-0/dbbench
du --total
cd /home/gopath/GoLevelBadgerRocksPebbleBenchmark
./clean.sh
cd /home/gopath/GoLevelBadgerRocksPebbleBenchmark/GoLevelDB
go run main.go
cd /tmp/leveldbtest-1/dbbench
du --total
cd /home/gopath/GoLevelBadgerRocksPebbleBenchmark
./clean.sh
cd /home/gopath/GoLevelBadgerRocksPebbleBenchmark/RocksDB
CGO_CFLAGS="-I/usr/local/rocksdb/include" CGO_LDFLAGS="-L/usr/local/rocksdb -lrocksdb -lstdc++ -lm -lz -lbz2 -lsnappy -llz4 -lzstd" go run main.go
cd /tmp/rocksdbtest-0/dbbench
du --total
cd /home/gopath/GoLevelBadgerRocksPebbleBenchmark
./clean.sh