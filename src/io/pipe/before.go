package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
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
	buf := bytes.NewBuffer(nil)
	json.NewEncoder(buf).Encode(msg) // 将msg序列化为byte,存到buf中
	resp, err := cli.Post("http://localhost:9999/json", "application/json", buf)

	if err != nil {
		log.Fatalln(err)
	}

	body := resp.Body
	defer body.Close()

	if body_bytes, err := ioutil.ReadAll(body); err == nil {
		log.Println("response:", string(body_bytes))
	} else {
		log.Fatalln(err)
	}
}
