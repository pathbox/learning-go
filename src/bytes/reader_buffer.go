package main

import (
	"bytes"
	"encoding/base64"
	"io"
	"os"
)

func main() {
  buf := bytes.NewBufferString("R29waGVycyBydWxlIQ==")
  dec := base64.NewDecoder(base64.StdEncoding, buf)
  io.Copy(os.Stdout, dec)
}
