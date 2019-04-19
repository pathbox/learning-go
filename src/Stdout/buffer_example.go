package main

import (
	"bytes"
	"io"
	"os"
)

var Out = os.Stdout

type Writer struct {
	Out io.Writer

	buf bytes.Buffer
}

func main() {
	w := &Writer{
		Out: Out,
	}

	for i := 0; i < 30; i++ {
		w.buf.WriteString("Hello World!\n")
	}

	w.Flush()
}

func (w *Writer) Flush() {
	w.Out.Write(w.buf.Bytes())
}
