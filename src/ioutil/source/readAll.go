import (
	"bytes"
	"io"
)

// ReadAll reads from r until an error or EOF and returns the data it read.
// A successful call returns err == nil, not err == EOF. Because ReadAll is
// defined to read from src until EOF, it does not treat an EOF from Read
// as an error to be reported.
func ReadAll(r io.Reader) ([]byte, error) {
	return readAll(r, bytes.MinRead)
}

func readAll(r io.Reader, capacity int64) (b []byte, err error) {
	var buf bytes.Buffer

	defer func() {
		e := recover()
		if e == nil {
			return
		}
		if panicErr, ok := e.(error); ok && panicErr == bytes.ErrTooLarge {
			err = panicErr
		} else {
			panic(e)
		}
	}()

	if int64(int(capacity)) == capacity {
		buf.Grow(int(capacity))
	}
	_, err = buf.ReadFrom(r) // 底层是 bytes.Buffer.ReadFrom 缓冲读
	return buf.Bytes(), err
}