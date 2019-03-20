package mmap

import (
	"errors"
	"os"
	"syscall"
)

var (
	// ErrUnmappedMemory is returned when a function is called on unmapped memory
	ErrUnmappedMemory = errors.New("unmapped memory")
	// ErrIndexOutOfBound is returned when given offset lies beyond the mapped region
	ErrIndexOutOfBound = errors.New("offset out of mapped region")
)

// Mmap provides abstraction around a memory mapped file
type Mmap struct {
	data   []byte
	length int64
}

// NewSharedFileMmap maps a file into memory starting at a given offset, for given length.
// For documentation regarding prot, see documentation for syscall package.
// possible cases:
//    case 1 => if   file size > memory region (offset + length)
//              then all the mapped memory is accessible
//    case 2 => if   file size <= memory region (offset + length)
//              then from offset to file size memory region is accessible
func NewSharedFileMmap(f *os.File, offset int64, length int, prot int) (IMmap, error) {
	data, err := syscall.Mmap(int(f.Fd()), offset, length, prot, syscall.MAP_SHARED)
	if err != nil {
		return nil, err
	}

	return &Mmap{
		data:   data,
		length: int64(length),
	}, nil
}

// Unmap unmaps the memory mapped file. An error will be returned
// if any of the functions are called on Mmap after calling Unmap
func (m *Mmap) Unmap() error {
	err := syscall.Munmap(m.data)
	m.data = nil
	return err
}
