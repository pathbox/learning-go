package main

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"
)

func main() {

}

func TestFile_SequenceWrite(t *testing.T) {
	f, _ := os.Create("./output.txt")
	defer f.Close()
	t1 := time.Now()
	for i := 0; i < 2000000; i++ {
		_, err := f.Write([]byte("0x6766c3279a7b32e52e89b24d203dd311aaf3019f9dd182f0128d8f12ab4490c2"))
		if err != nil {
			t.Fatalf("put failed: %v", err)
		}
	}
	t2 := time.Now()
	fmt.Println("spend time:", t2.Sub(t1))
}

func TestFile_ConcurrentWrite(t *testing.T) {
	f1, _ := os.Create("./output1.txt")
	f2, _ := os.Create("./output2.txt")
	wg := sync.WaitGroup{}
	wg.Add(2)
	t1 := time.Now()
	go func() {
		for i := 0; i < 1000000; i++ {
			_, err := f1.Write([]byte("0x6766c3279a7b32e52e89b24d203dd311aaf3019f9dd182f0128d8f12ab4490c2"))
			if err != nil {
				t.Fatalf("put failed: %v", err)
			}
		}
		wg.Done()
	}()
	go func() {
		for i := 1000000; i < 2000000; i++ {
			_, err := f2.Write([]byte("0x6766c3279a7b32e52e89b24d203dd311aaf3019f9dd182f0128d8f12ab4490c2"))
			if err != nil {
				t.Fatalf("put failed: %v", err)
			}
		}
		wg.Done()
	}()
	wg.Wait()
	t2 := time.Now()
	fmt.Println("spend time:", t2.Sub(t1))
}

/*
经测试，写入一个66字节的字符串200万次 到一个文件里（约134.1MB），如果顺序执行，需要耗时12秒左右。但如果开两个协程并发写入两个不同的文件，耗时在8秒左右。如果开两个协程以上并发写入多个文件，不会再有性能的提高。当然，这个速度和固态硬盘的读写速度是差很远的，如果想提高写入速度，可以在内存里先构造出需要写入文件的内容(先写入buffer，然后再批量的flush一次到磁盘)，再写入，这样能极大提高写入速度。如果使用这种方式写入文件数据，以上述的环境只需要30-40µs即可写入完成
*/
