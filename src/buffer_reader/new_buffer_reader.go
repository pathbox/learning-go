package buffer_reader

import (
	"bufio"
	"bytes"
)

func New() *bufio.Reader {
	var bs []byte

	reader := bytes.NewReader(bs)
	bufReader := bufio.NewReader(reader)
	return bufReader
}

/*
// bytes.NewReader() NewReader returns a new Reader reading from b.
func NewReader(b []byte) *Reader { return &Reader{b, 0, -1} }

// bufio.NewReader() NewReader returns a new Reader whose buffer has the default size.
func NewReader(rd io.Reader) *Reader {
	return NewReaderSize(rd, defaultBufSize) // 4096
}

func NewReaderSize(rd io.Reader, size int) *Reader {
	b, ok := rd.(*Reader)
	if ok && len(b.buf) >= size {
		return b
	}

	if size < minReadBufferSize {
		size = minReadBufferSize // 16
	}
	r := new(Reader)
	r.reset(make([]byte, size), rd)
	return r
}

func (b *Reader) reset(buf []byte, r io.Reader) {
	*b = Reader{
		buf:          buf,
		rd:           r,
		lastByteZ:    -1,
		lastRuneSize: -1,
	}
}

*/
