project_path=$GOPATH/src/github.com/SigelKoo
data_path=${project_path}/GoLevelBadgerRocksPebbleBenchmark2/data
echo "data_path: ${data_path}"

rm -rf ${data_path}/badgertest-0/dbbench/*
rm -rf ${data_path}/leveldbtest-1/dbbench/*
rm -rf ${data_path}/pebbledbtest-0/dbbench/*
rm -rf ${data_path}/rocksdbtest-0/dbbench/*