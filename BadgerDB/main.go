package main

import (
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"time"
	"unsafe"

	"github.com/dgraph-io/badger/v3"
)

var (
	keySz        int  = 64
	valueSz      int  = 1024
	dataCntRange int  = 10000
	batchCnt     int  = 1000
	sync         bool = true
)

func main() {

	badgerTimes := float64(0)
	badgerSpeeds := make([]float64, 0, dataCntRange)
	for i := 1; i <= dataCntRange; i++ {
		bs, bt := bench_test()
		badgerTimes += bt
		badgerSpeeds = append(badgerSpeeds, bs)
	}

	for i := 0; i < len(badgerSpeeds); i++ {
		fmt.Printf("total: %d, badgerTime: %f μs/op\n",
			(i+1)*batchCnt, badgerSpeeds[i])
	}
	fmt.Printf("Usage Time: %f μs\n", badgerTimes)
	fmt.Println()
}

func bench_test() (badgerSpeed, badgerTime float64) {
	db := BadgerConn()
	btotalWriteTime := float64(0)
	keyList := [][]byte{}
	valueList := [][]byte{}
	for j := 0; j < batchCnt; j++ {
		keyList = append(keyList, RandStr(keySz))
		valueList = append(valueList, RandStr(valueSz))
	}
	//pstart := time.Now()
	wstart := time.Now()
	wb := db.NewWriteBatch()
	for j := 0; j < batchCnt; j++ {
		wb.Set(keyList[j], valueList[j])
	}
	//pend := time.Since(pstart)
	keyList = nil
	valueList = nil
	//btotalWriteTime += float64(pend.Microseconds())
	//wstart := time.Now()
	err := wb.Flush()
	wend := time.Since(wstart)
	if err != nil {
		log.Fatal(err)
	}
	btotalWriteTime += float64(wend.Microseconds())
	btotalWriteSpeed := btotalWriteTime / float64(batchCnt)
	db.Close()
	return btotalWriteSpeed, btotalWriteTime
}

func BadgerConn() *badger.DB {
	opts := badger.DefaultOptions("../data/badgertest-0/dbbench")
	opts.SyncWrites = sync
	db, err := badger.Open(opts)
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
