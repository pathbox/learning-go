package main

import (
	"bytes"
	"fmt"
	"os"
)

func main() {
  var b bytes.Buffer
  b.Write([]byte("Hello "))
  fmt.Println(&b)
  fmt.Fprintf(&b, "World!\n")
  fmt.Println(&b)
  b.WriteTo(os.Stdout)
}
