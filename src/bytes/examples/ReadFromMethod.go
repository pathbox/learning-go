import "io"

// ReadFrom reads data from r(io.Reader) until EOF and appends it to the buffer, growing
// the buffer as needed. The return value n is the number of bytes read. Any
// error except io.EOF encountered during the read is also returned. If the
// buffer becomes too large, ReadFrom will panic with ErrTooLarge
// 循环读取，每次读512长度的字节数据。不会消耗内存
func (b *Buffer) ReadFrom(r io.Reader) (n int64, err error) {
	b.lastRead = opInvalid
	for {
		i := b.grow(MinRead)                // 每次读取之前，将buf扩容 512 长度的cap，用于下一次读取数据，每次最多读取MinRead
		m, e := r.Read(b.buf[i:cap(b.buf)]) // data from reader to => buf, len(data) is m
		if m < 0 {
			panic(errNegativeRead)
		}
		b.buf = b.buf[:i+m] // 修改buf
		n += int64(m)       // total read len(data)
		if e == io.EOF {
			return n, nil // e is EOF, so return nil explicitly
		}
		if e != nil {
			return n, e
		}
	}
}