ulimit -a

project_path=$GOPATH/src/github.com/SigelKoo
data_path=${project_path}/GoLevelBadgerRocksPebbleBenchmark2/data
echo "data_path: ${data_path}"

cd ${project_path}/GoLevelBadgerRocksPebbleBenchmark2/BadgerDB
go run main.go
# cd ${data_path}/badgertest-0/dbbench 
#du --total
cd ${project_path}/GoLevelBadgerRocksPebbleBenchmark2
#./clean.sh
cd ${project_path}/GoLevelBadgerRocksPebbleBenchmark2/PebbleDB
go run main.go
# cd ${data_path}/pebbledbtest-0/dbbench
#du --total
cd ${project_path}/GoLevelBadgerRocksPebbleBenchmark2
#./clean.sh
cd ${project_path}/GoLevelBadgerRocksPebbleBenchmark2/GoLevelDB
go run main.go
# cd ${data_path}/leveldbtest-1/dbbench
#du --total
cd ${project_path}/GoLevelBadgerRocksPebbleBenchmark2
#./clean.sh
cd ${project_path}/GoLevelBadgerRocksPebbleBenchmark2/RocksDB
CGO_CFLAGS="-I/usr/local/rocksdb/include" CGO_LDFLAGS="-L/usr/local/rocksdb -lrocksdb -lstdc++ -lm -lz -lbz2 -lsnappy -llz4 -lzstd" go run main.go
# cd ${data_path}/rocksdbtest-0/dbbench
#du --total
cd ${project_path}/GoLevelBadgerRocksPebbleBenchmark2
#./clean.sh