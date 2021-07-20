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
	valueSz      int  = 10
	dataCntRange int  = 10000
	skip         int  = 1
	batchCnt     int  = 1000
	sync         bool = false
)

func main() {

	pebbleTimes := make([]float64, 0, dataCntRange)
	for i := 1; i <= dataCntRange; i++ {
		bt := bench_test(i * skip)
		pebbleTimes = append(pebbleTimes, bt)
	}

	for i := 0; i < len(pebbleTimes); i++ {
		fmt.Printf("total: %d, pebbleTimes: %f μs/op",
			(i+1)*batchCnt*skip, pebbleTimes[i])
	}
}

func bench_test(dataCnt int) (pebbleTime float64) {
	total := dataCnt * batchCnt

	db := PebbleConn()
	fmt.Println("Pebble:")
	pstart := time.Now()
	ptotalWriteTime := float64(0)
	for i := 0; i < dataCnt; i++ {
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
		ptotalWriteTime = ptotalWriteTime + float64(pend.Microseconds())
		wstart := time.Now()
		batch.Commit(&pebble.WriteOptions{Sync: sync})
		wend := time.Since(wstart)
		//fmt.Printf("badger write %d st data\n", i)
		ptotalWriteTime = ptotalWriteTime + float64(wend.Microseconds())
	}
	ptotalWriteTime = ptotalWriteTime / float64(total)
	//fmt.Printf("Total write time: %f μs/op\n", ptotalWriteTime)
	fmt.Println("Total time: ", time.Since(pstart))
	db.Close()
	return ptotalWriteTime
}

func PebbleConn() *pebble.DB {
	db, err := pebble.Open("/tmp/pebbledbtest-0/dbbench", &pebble.Options{})
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
