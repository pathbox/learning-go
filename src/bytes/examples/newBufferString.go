package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

func main() {
	// A Buffer can turn a string or a []byte into an io.Reader.
	buf := bytes.NewBufferString("R29waGVycyBydWxlIQ==")
	dec := base64.NewDecoder(base64.StdEncoding, buf)
	fmt.Println(buf.Bytes())
	io.Copy(os.Stdout, dec)
}
