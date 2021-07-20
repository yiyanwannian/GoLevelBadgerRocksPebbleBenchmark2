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
	valueSz      int  = 10
	dataCntRange int  = 10000
	skip         int  = 1
	batchCnt     int  = 1000
	sync         bool = false
)

func main() {

	rocksDBTimes := make([]float64, 0, dataCntRange)
	for i := 1; i <= dataCntRange; i++ {
		bt := bench_test(i * skip)
		rocksDBTimes = append(rocksDBTimes, bt)
	}

	for i := 0; i < len(rocksDBTimes); i++ {
		fmt.Printf("total: %d, rocksDBTimes: %f μs/op",
			(i+1)*batchCnt*skip, rocksDBTimes[i])
	}
}

func bench_test(dataCnt int) (rocksDBTime float64) {
	total := dataCnt * batchCnt

	db := RocksDbConn()
	wo := gorocksdb.NewDefaultWriteOptions()
	wo.SetSync(sync)
	fmt.Println("RocksDB:")
	rstart := time.Now()
	rtotalWriteTime := float64(0)
	for i := 0; i < dataCnt; i++ {
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
		rtotalWriteTime = rtotalWriteTime + float64(pend.Microseconds())
		wstart := time.Now()
		err := db.Write(wo, wb)
		wend := time.Since(wstart)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Printf("badger write %d st data\n", i)
		rtotalWriteTime = rtotalWriteTime + float64(wend.Microseconds())
	}
	rtotalWriteTime = rtotalWriteTime / float64(total)
	//fmt.Printf("Total write time: %f μs/op\n", rtotalWriteTime)
	fmt.Println("Total time: ", time.Since(rstart))
	db.Close()
	return rtotalWriteTime
}

func RocksDbConn() *gorocksdb.DB {
	opts := gorocksdb.NewDefaultOptions()
	opts.SetCreateIfMissing(true)
	db, err := gorocksdb.OpenDb(opts, "/tmp/rocksdbtest-0/dbbench")
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
