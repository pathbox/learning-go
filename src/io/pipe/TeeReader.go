package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

func main() {
	pr, pw := io.Pipe()

	wg := sync.WaitGroup{}
	wg.Add(2)

	f, err := os.Open("./fruit.txt")
	if err != nil {
		log.Fatal(err)
	}

	tr := io.TeeReader(f, pw)  // 将　ｆ　通过ｐｗ进行写入操作

	go func() {
		defer wg.Done()
		defer pw.Close()

		// get data from the TeeReader, which feeds the PipeReader through the PipeWriter
		_, err := http.Post("https://example.com", "text/html", tr)
		if err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		defer wg.Done()

		if _, err := io.Copy(os.Stdout, pr); err != nil {　// 将ｐｗ写入的数据　通ｐｒ读取到os.Stdout
			log.Fatal(err)
		}
	}()

	wg.Wait()
}
