import "io"

func (b *Buffer) WriteTo(w io.Writer) (n int64, err error) {
	b.lastRead = opInvalid
	if nBytes := b.Len(); nBytes > 0 { // 将b中的数据全部写入到w
		m, e := w.Write(b.buf[b.off:])
		if m > nBytes {
			panic("bytes.Buffer.WriteTo: invalid Write count")
		}
		b.off += m // off 就是byte指针偏移量
		n = int64(m)
		if e != nil {
			return n, e
		}

		if m != nBytes {
			return n, io.ErrShortWrite
		}
	}

	b.Reset() // b重新置空
	return n, nil
}