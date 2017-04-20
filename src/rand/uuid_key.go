package main

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"time"
)

func main() {
	id := newId()
	fmt.Println(id)
}

func newId() string {
	hash := fmt.Sprintf("%s", time.Now())
	buf := bytes.NewBuffer(nil)
	sum := md5.Sum([]byte(hash))
	encoder := base64.NewEncoder(base64.URLEncoding, buf)
	encoder.Write(sum[:])
	encoder.Close()
	return buf.String()[:20]
}
