package markov

import (
	"bytes"
	"io"
	"log"
)

func (c *Chain) Write(p []byte) (int, error) {
	log.Printf("chain received: %s", p)
	c.Build(bytes.NewReader(p))
	return len(p), nil
}

func (c *Chain) SpyOn(r io.Reader) io.Reader {
	pr, pw := io.Pipe()
	go func() {
		_, err := io.Copy(io.MultiWriter(pw, c), r)
		pw.CloseWithError(err)
	}()
	return pr
}
