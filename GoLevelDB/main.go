package main

import (
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"time"
	"unsafe"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

var (
	keySz        int  = 64
	valueSz      int  = 1024
	dataCntRange int  = 10
	batchCnt     int  = 1000
	sync         bool = true
)

func main() {
	leveldbTimes := float64(0)
	leveldbSpeeds := make([]float64, 0, dataCntRange)
	for i := 1; i <= dataCntRange; i++ {
		bs, bt := bench_test()
		leveldbTimes += bt
		leveldbSpeeds = append(leveldbSpeeds, bs)
	}

	for i := 0; i < dataCntRange; i++ {
		fmt.Printf("total: %d, levelDBTime: %f μs/op\n",
			(i+1)*batchCnt, leveldbSpeeds[i])
	}
	fmt.Printf("Usage Time: %f μs\n", leveldbTimes)
	fmt.Println()
}

func bench_test() (leveldbSpeed, levelDBTime float64) {
	db := GoLevelDbConn()
	ltotalWriteTime := float64(0)
	keyList := [][]byte{}
	valueList := [][]byte{}
	for j := 0; j < batchCnt; j++ {
		keyList = append(keyList, RandStr(keySz))
		valueList = append(valueList, RandStr(valueSz))
	}
	wstart := time.Now()
	//pstart := time.Now()
	batch := new(leveldb.Batch)
	for j := 0; j < batchCnt; j++ {
		batch.Put(keyList[j], valueList[j])
	}
	//pend := time.Since(pstart)
	keyList = nil
	valueList = nil
	//ltotalWriteTime += float64(pend.Microseconds())
	//wstart := time.Now()
	err := db.Write(batch, &opt.WriteOptions{Sync: sync})
	wend := time.Since(wstart)
	if err != nil {
		log.Fatal(err)
	}
	ltotalWriteTime += float64(wend.Microseconds())
	ltotalWriteSpeed := ltotalWriteTime / float64(batchCnt)
	db.Close()
	return ltotalWriteSpeed, ltotalWriteTime
}

func GoLevelDbConn() *leveldb.DB {
	opts := &opt.Options{NoSync: !sync}
	db, err := leveldb.OpenFile("../data/leveldbtest-1/dbbench", opts)
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
