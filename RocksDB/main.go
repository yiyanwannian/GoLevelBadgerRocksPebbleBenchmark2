package main

import (
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"time"
	"unsafe"

	"github.com/tecbot/gorocksdb"
)

var (
	keySz        int  = 64
	valueSz      int  = 1024
	dataCntRange int  = 10
	batchCnt     int  = 1000
	sync         bool = true
)

func main() {
	rocksDBTimes := float64(0)
	rocksDBSpeeds := make([]float64, 0, dataCntRange)
	for i := 1; i <= dataCntRange; i++ {
		bs, bt := bench_test()
		rocksDBTimes += bt
		rocksDBSpeeds = append(rocksDBSpeeds, bs)
	}

	for i := 0; i < dataCntRange; i++ {
		fmt.Printf("total: %d, rocksDBTimes: %f μs/op\n",
			(i+1)*batchCnt, rocksDBSpeeds[i])
	}
	fmt.Println()
	fmt.Printf("Usage Time: %f μs\n", rocksDBTimes)
}

func bench_test() (rocksDBSpeed, rocksDBTime float64) {
	db := RocksDbConn()
	wo := gorocksdb.NewDefaultWriteOptions()
	wo.SetSync(sync)
	rtotalWriteTime := float64(0)
	keyList := [][]byte{}
	valueList := [][]byte{}
	for j := 0; j < batchCnt; j++ {
		keyList = append(keyList, RandStr(keySz))
		valueList = append(valueList, RandStr(valueSz))
	}
	pstart := time.Now()
	wb := gorocksdb.NewWriteBatch()
	for j := 0; j < batchCnt; j++ {
		wb.Put(keyList[j], valueList[j])
	}
	pend := time.Since(pstart)
	keyList = nil
	valueList = nil
	rtotalWriteTime += float64(pend.Microseconds())
	wstart := time.Now()
	err := db.Write(wo, wb)
	wend := time.Since(wstart)
	if err != nil {
		log.Fatal(err)
	}
	rtotalWriteTime += float64(wend.Microseconds())
	rtotalWriteSpeed := rtotalWriteTime / float64(batchCnt)
	db.Close()
	return rtotalWriteSpeed, rtotalWriteTime
}

func RocksDbConn() *gorocksdb.DB {
	opts := gorocksdb.NewDefaultOptions()
	opts.SetCreateIfMissing(true)
	opts.SetMaxBackgroundCompactions(4)
	opts.SetMaxBackgroundFlushes(2)
	db, err := gorocksdb.OpenDb(opts, "../data/rocksdbtest-0/dbbench")
	if err != nil {
		log.Fatal("err")
	}
	return db
}

func RandStr(length int) []byte {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := String2Bytes(str)
	result := []byte{}
	rand.Seed(time.Now().UnixNano() + int64(rand.Intn(100)))
	for i := 0; i < length; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return result
}

func String2Bytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}
