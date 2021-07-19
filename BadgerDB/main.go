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
	skip         int  = 1
	batchCnt     int  = 1000
	sync         bool = false
)

func main() {

	badgerTimes := make([]float64, 0, dataCntRange)
	for i := 1; i <= dataCntRange; i++ {
		bt := bench_test(i * skip)
		badgerTimes = append(badgerTimes, bt)
	}

	for i := 0; i < len(badgerTimes); i++ {
		fmt.Printf("total: %d, badgerTime: %f μs/op",
			(i+1)*batchCnt*skip, badgerTimes[i])
	}
}

func bench_test(dataCnt int) (badgerTime float64) {
	total := dataCnt * batchCnt

	db := BadgerConn()
	fmt.Println("Badger:")
	bstart := time.Now()
	btotalWriteTime := float64(0)
	for i := 0; i < dataCnt; i++ {
		wb := db.NewWriteBatch()
		for j := 0; j < batchCnt; j++ {
			key := RandStr(keySz)
			value := RandStr(valueSz)
			pstart := time.Now()
			wb.Set(key, value)
			pend := time.Since(pstart)
			btotalWriteTime = btotalWriteTime + float64(pend.Microseconds())
		}
		wstart := time.Now()
		err := wb.Flush()
		wend := time.Since(wstart)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Printf("badger write %d st data\n", i)
		btotalWriteTime = btotalWriteTime + float64(wend.Microseconds())
	}
	btotalWriteTime = btotalWriteTime / float64(total)
	//fmt.Printf("Total write time: %f μs/op\n", btotalWriteTime)
	fmt.Println("Total time: ", time.Since(bstart))
	db.Close()
	return btotalWriteTime
}

func BadgerConn() *badger.DB {
	opts := badger.DefaultOptions("/tmp/badgertest-0/dbbench")
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
