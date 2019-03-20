package mmap

import "io"

type IMmap interface {
	io.ReaderAt
	io.WriteAt

	Lock() error
	Unlock() error
	Advise(advice int) error
	ReadUint64At(offset int64) uint64
	WriteUint64At(num uint64, offset int64)
	Flush(flags int) error
	Unmap() error
}
