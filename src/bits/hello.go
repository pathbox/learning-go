package main

import (
	"fmt"
	"os"

	"github.com/nareix/bits"
)

func main() {
  writer, _ := os.Create("test.bin")
  w := &bits.Writer{
    W: writer,
  }
  err := w.WriteBits(0x22, 4)
  if err != nil{
    panic(err)
  }
  err = w.WriteBits64(0x22, 4)
  n, err := w.Write([]byte{0x11, 0x22, 0x33, 0x44})
  if err != nil {
    panic(err)
  }
  fmt.Printf("write %d bytes\n", n)
  err = w.FlushBits()

  reader, _ := os.Open("test.bin")
  r := &bits.Reader{
    R: reader,
  }
  u32, _ := r.ReadBits(4)
  fmt.Printf("u32: %v\n", u32)
  u64, _ := r.ReadBits64(4)
  fmt.Printf("u64: %v\n", u64)
  p := make([]byte, 4)
  n, err = r.Read(p)
  if err != nil {
    panic(err)
  }
  fmt.Printf("byte %d: %v\n", n, p)
}
