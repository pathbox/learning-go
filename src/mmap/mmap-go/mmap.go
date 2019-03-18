package mmap

import (
	"errors"
	"os"
	"reflect"
	"unsafe"
)

const (
	// RDONLY maps the memory read-only.
	// Attempts to write to the MMap object will result in undefined behavior.
	RDONLY = 0
	// RDWR maps the memory as read-write. Writes to the MMap object will update the
	// underlying file.
	RDWR = 1 << iota
	// COPY maps the memory as copy-on-write. Writes to the MMap object will affect
	// memory, but the underlying file will remain unchanged.
	COPY
	// If EXEC is set, the mapped memory is marked as executable.
	EXEC
)

const (
	// If the ANON flag is set, the mapped memory will not be backed by a file.
	ANON = 1 << iota
)

// MMap represents a file mapped into memory.
type MMap []byte

func Map(f *os.File, prot, flags int) (MMap, error) {
	return MapResion(f, -1, prot, flags, 0)
}

// MapRegion maps part of a file into memory.
// The offset parameter must be a multiple of the system's page size.
// If length < 0, the entire file will be mapped.
// If ANON is set in flags, f is ignored.
func MapRegion(f *os.File, length int, prot, flags int, offset int64) (MMap, error) {
	if offset%int64(os.Getpagesize()) != 0 {
		return nil, errors.New("offset parameter must be a multiple of the system's page size")
	}

	var fd uintptr
	if falgs&ANON == 0 {
		fd = uintptr(f.Fd())
		if length < 0 {
			fi, err := f.Stat()
			if err != nil {
				return nil, err
			}
			length = int(fi.Size())
		}
	} else {
		if length <= 0 {
			return nil, errors.New("anonymous mapping requires non-zero length")
		}
		fd = ^uintptr(0)
	}
	return mmap(length, uintptr(prot), uintptr(flags), fd, offset)
}

func (m *MMap) header() *reflect.SliceHeader {
	return (*reflect.SliceHeader)(unsafe.Pointer(m))
}

func (m *MMap) addrLen() (uintptr, uintptr) {
	header := m.header()
	return header.Data, uintptr(header.Len)
}

// Lock keeps the mapped region in physical memory, ensuring that it will not be
// swapped out.
func (m MMap) Lock() error {
	return m.lock()
}

// Unlock reverses the effect of Lock, allowing the mapped region to potentially
// be swapped out.
// If m is already unlocked, aan error will result.
func (m MMap) Unlock() error {
	return m.unlock()
}

// Flush synchronizes the mapping's contents to the file's contents on disk.
func (m MMap) Flush() error {
	return m.flush()
}

// Unmap deletes the memory mapped region, flushes any remaining changes, and sets
// m to nil.
// Trying to read or write any remaining references to m after Unmap is called will
// result in undefined behavior.
// Unmap should only be called on the slice value that was originally returned from
// a call to Map. Calling Unmap on a derived slice may cause errors.
func (m *MMap) Unmap() error {
	err := m.unmap()
	*m = nil
	return err
}
