package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func init() {
	log.SetFlags(log.Lshortfile)
}
func main() {
	cli := http.Client{}

	msg := struct {
		Name, Addr string
		Price      float64
	}{
		Name:  "hello",
		Addr:  "beijing",
		Price: 123.56,
	}
	r, w := io.Pipe()
	// 注意这边的逻辑！！
	go func() {
		defer func() {
			time.Sleep(time.Second * 2)
			log.Println("encode完成")
			// 只有这里关闭了，Post方法才会返回
			w.Close()
		}()
		log.Println("管道准备输出")
		// 只有Post开始读取数据，这里才开始encode，并传输
		err := json.NewEncoder(w).Encode(msg) // 在一个新的goroutine中,进行write操作
		log.Println("管道输出数据完毕")
		if err != nil {
			log.Fatalln("encode json failed:", err)
		}
	}()
	time.Sleep(time.Second * 1)
	log.Println("开始从管道读取数据")
	resp, err := cli.Post("http://localhost:9999/json", "application/json", r) // 将write操作写入的byte读取出来

	if err != nil {
		log.Fatalln(err)
	}
	log.Println("POST传输完成")

	body := resp.Body
	defer body.Close()

	if body_bytes, err := ioutil.ReadAll(body); err == nil {
		log.Println("response:", string(body_bytes))
	} else {
		log.Fatalln(err)
	}
}
