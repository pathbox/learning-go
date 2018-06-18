import "io"

type Writer struct {
	err error
	buf []byte
	n   int
	wr  io.Writer
}

func NewWriterSize(w io.Writer, size int) *Writer {
	b, ok := w.(*Writer)
	if ok && len(b.buf) >= size {
		return b
	}
	if size <= 0 {
		size = defaultBufSize
	}
	return &Writer{
		buf: make([]byte, size),
		wr:  w,
	}
}

func (b *Writer) Reset(w io.Writer {
	b.err = nil
	b.n = 0
	b.wr = w
})

func (b *Writer) Flush() error {
	if b.err != nil {
		return b.err
	}
	if b.n == 0 {
		return nil
	}
	n, err := b.wr.Write(b.buf[0:b.n]) // write all buf bytes to io.writer
	if n < b.n && err == nil {
		err = io.ErrShortWrite
	}
	if err != nil {
		if n > 0 && n < b.n {
			copy(b.buf[0:b.n-n],b.buf[n:b.n])
		}
		b.n -= n
		b.err = err
		return err
	}
	b.n = 0
	return nil
}

func (b *Writer) Write(p []byte) (nn int, err error) {
	for len(p) > b.Avaiable() && b.err == nil {
		var n int
		if b.Buffered() == 0 {
			n, b.err = b.wr.Write(p)
		} else {
			n = copy(b.buf[b.n:], p)
			b.n += n
			b.Flush()
		}
		nn += n
		p = p[n:]
	}
	if b.err != nil {
		return nn, b.err
	}
	n := copy(b.buf[b.n:], p)
	b.n += n
	nn += n
	return nn, nil
}

// bufio 包 实际上是实现了一个缓冲层buf, 将数据先存到其缓冲buf []byte 中，之后再write 到其他io.Writer中