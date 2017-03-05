type Reader interface {
  Read(p []byte) (n int, err error)
}

type Writer interface {
  Write(p []byte) (n int, err error)
}

type Seeker interface {
  Seek(offset int64, whence int) (int64, error)
}

type ReaderAt interface {
  ReadAt(p []byte, off int64) (n int, err error)
}

type WriterAt interface {
  WriteAt(p []byte, off int64) (n int, err error)
}

type Closer interface {
  Close() error
}


// go mix-in

type ReadCloser interface {
  Reader
  Closer
}

type ReadSeeker interface {
  Reader
  Seeker
}

type WriteCloser interface {
  Writer
  Seeker
}

type ReadWriter interface {
  Reader
  Writer
}

type ReadWriteCloser interface {
  Reader
  Writer
  Closer
}

type ReadWriteSeeker interface {
  Reader
  Writer
  Seeker
}

type WriterTo interface {
    WriteTo(w Writer) (n int64, err error)
}
type ReaderFrom interface {
    ReadFrom(r Reader) (n int64, err error)
}

type ByteReader interface {
  ReadByte() (c byte, err error)
}

type ByteScanner interface {
  ByteReader
  UnReadByte() error
}

type ByteWriter interface {
  WriteByte(c byte) error
}

type RunReader interface {
  ReadRune() (r rune, size int, err error)
}

type RuneScanner interface {
  RuneReader
  UnreadRune() error
}