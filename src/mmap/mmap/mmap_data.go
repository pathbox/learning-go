package mmap

import (
	"encoding/binary"
	"syscall"
	"unsafe"
)

// ReadAt copies data to dest slice from mapped region starting at
// given offset and returns number of bytes copied to the dest slice.
// There are two possibilities -
//   Case 1: len(dest) >= (len(m.data) - offset)
//        => copies (len(m.data) - offset) bytes to dest from mapped region
//   Case 2: len(dest) < (len(m.data) - offset)
//        => copies len(dest) bytes to dest from mapped region
// err is always nil, hence, can be ignored
func (m *Mmap) ReadAt(dest []byte, offset int64) (int, error) {
	if m.data == nil {
		panic(ErrUnmappedMemory)
	} else if offset >= m.length || offset < 0 {
		panic(ErrIndexOutOfBound)
	}

	return copy(dest, m.data[offset:]), nil
}

// WriteAt copies data to mapped region from the src slice starting at
// given offset and returns number of bytes copied to the mapped region.
// There are two possibilities -
//  Case 1: len(src) >= (len(m.data) - offset)
//      => copies (len(m.data) - offset) bytes to the mapped region from src
//  Case 2: len(src) < (len(m.data) - offset)
//      => copies len(src) bytes to the mapped region from src
// err is always nil, hence, can be ignored
func (m *Mmap) WriteAt(src []byte, offset int64) (int, error) {
	if m.data == nil {
		panic(ErrUnmappedMemory)
	} else if offset >= m.length || offset < 0 {
		panic(ErrIndexOutOfBound)
	}

	return copy(m.data[offset:], src), nil
}

// ReadUint64At reads uint64 from offset
func (m *Mmap) ReadUint64At(offset int64) uint64 {
	if m.data == nil {
		panic(ErrUnmappedMemory)
	} else if offset+8 > m.length || offset < 0 {
		panic(ErrIndexOutOfBound)
	}

	return binary.LittleEndian.Uint64(m.data[offset : offset+8])
}

// WriteUint64At writes num at offset
func (m *Mmap) WriteUint64At(num uint64, offset int64) {
	if m.data == nil {
		panic(ErrUnmappedMemory)
	} else if offset+8 > m.length || offset < 0 {
		panic(ErrIndexOutOfBound)
	}

	binary.LittleEndian.PutUint64(m.data[offset:offset+8], num)
}

// Flush flushes the memory mapped region to disk
func (m *Mmap) Flush(flags int) error {
	_, _, err := syscall.Syscall(syscall.SYS_MSYNC,
		uintptr(unsafe.Pointer(&m.data[0])), uintptr(len(m.data)), uintptr(flags))
	if err != 0 {
		return err
	}

	return nil
}
