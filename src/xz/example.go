package main

import (
	"bytes"
	"io"
	"log"
	"os"

	xz "github.com/ulikunitz/xz"
)

func main() {
	const text = "The quick brown fox jumps over the lazy dog.\n"
	var buf bytes.Buffer
	w, err := xz.NewWriter(&buf)
	if err != nil {
		log.Fatalf("xz.NewWriter error %s", err)
	}
	if _, err := io.WriteString(w, text); err != nil {
		log.Fatalf("WriteString error %s", err)
	}
	if err := w.Close(); err != nil {
		log.Fatalf("w.Close error %s", err)
	}

	// pass buf data by Internet

	r, err := xz.NewReader(&buf)
	if err != nil {
		log.Fatalf("NewReader error %s", err)
	}
	if _, err = io.Copy(os.Stdout, r); err != nil {
		log.Fatalf("io.Copy error %s", err)
	}
}
