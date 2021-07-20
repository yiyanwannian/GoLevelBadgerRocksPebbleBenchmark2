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
	valueSz      int  = 10
	dataCntRange int  = 10000
	skip         int  = 1
	batchCnt     int  = 1000
	sync         bool = false
)

func main() {

	leveldbTimes := make([]float64, 0, dataCntRange)
	for i := 1; i <= dataCntRange; i++ {
		bt := bench_test(i * skip)
		leveldbTimes = append(leveldbTimes, bt)
	}

	for i := 0; i < len(leveldbTimes); i++ {
		fmt.Printf("total: %d, badgerTime: %f μs/op",
			(i+1)*batchCnt*skip, leveldbTimes[i])
	}
}

func bench_test(dataCnt int) (leveldbTime float64) {
	total := dataCnt * batchCnt

	db := GoLevelDbConn()
	fmt.Println("LevelDB:")
	lstart := time.Now()
	ltotalWriteTime := float64(0)
	for i := 0; i < dataCnt; i++ {
		keyList := [][]byte{}
		valueList := [][]byte{}
		for j := 0; j < batchCnt; j++ {
			keyList = append(keyList, RandStr(keySz))
			valueList = append(valueList, RandStr(valueSz))
		}
		pstart := time.Now()
		batch := new(leveldb.Batch)
		for j := 0; j < batchCnt; j++ {
			batch.Put(keyList[j], valueList[j])
		}
		pend := time.Since(pstart)
		ltotalWriteTime = ltotalWriteTime + float64(pend.Microseconds())
		wstart := time.Now()
		err := db.Write(batch, &opt.WriteOptions{Sync: sync})
		wend := time.Since(wstart)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Printf("badger write %d st data\n", i)
		ltotalWriteTime = ltotalWriteTime + float64(wend.Microseconds())
	}
	ltotalWriteTime = ltotalWriteTime / float64(total)
	//fmt.Printf("Total write time: %f μs/op\n", ltotalWriteTime)
	fmt.Println("Total time: ", time.Since(lstart))
	db.Close()
	return ltotalWriteTime
}

func GoLevelDbConn() *leveldb.DB {
	opts := &opt.Options{NoSync: !sync}
	db, err := leveldb.OpenFile("/tmp/leveldbtest-1/dbbench", opts)
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
