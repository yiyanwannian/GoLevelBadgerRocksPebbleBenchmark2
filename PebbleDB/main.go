package main

import (
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"time"
	"unsafe"

	"github.com/cockroachdb/pebble"
)

var (
	keySz        int  = 64
	valueSz      int  = 1024
	dataCntRange int  = 10000
	batchCnt     int  = 1000
	sync         bool = true
)

func main() {
	pebbledbTimes := float64(0)
	pebbledbSpeeds := make([]float64, 0, dataCntRange)
	for i := 1; i <= dataCntRange; i++ {
		bs, bt := bench_test()
		pebbledbTimes += bt
		pebbledbSpeeds = append(pebbledbSpeeds, bs)
	}

	for i := 0; i < dataCntRange; i++ {
		fmt.Printf("total: %d, pebbleDBTime: %f μs/op",
			(i+1)*batchCnt, pebbledbSpeeds[i])
	}
	fmt.Println()
	fmt.Printf("Usage Time: %f μs\n", pebbledbTimes)
}

func bench_test() (pebbleDBSpeed, pebbleDBTime float64) {
	db := PebbleConn()
	ptotalWriteTime := float64(0)
	keyList := [][]byte{}
	valueList := [][]byte{}
	for j := 0; j < batchCnt; j++ {
		keyList = append(keyList, RandStr(keySz))
		valueList = append(valueList, RandStr(valueSz))
	}
	pstart := time.Now()
	batch := db.NewBatch()
	for j := 0; j < batchCnt; j++ {
		batch.Set(keyList[j], valueList[j], &pebble.WriteOptions{Sync: sync})
	}
	pend := time.Since(pstart)
	keyList = nil
	valueList = nil
	ptotalWriteTime += float64(pend.Microseconds())
	wstart := time.Now()
	batch.Commit(&pebble.WriteOptions{Sync: sync})
	wend := time.Since(wstart)
	ptotalWriteTime += float64(wend.Microseconds())
	ptotalWriteSpeed := ptotalWriteTime / float64(batchCnt)
	db.Close()
	return ptotalWriteSpeed, ptotalWriteTime
}

func PebbleConn() *pebble.DB {
	db, err := pebble.Open("/tmp/pebbledbtest-0/dbbench", &pebble.Options{MaxConcurrentCompactions:4})
	if err != nil {
		log.Fatal(err)
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
