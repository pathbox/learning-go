package mmap

import (
	"os"
	"syscall"
)

// Open will mmap a file to a byte slice of data.
func Open(path string, writable bool) (data []byte, err error) {
	flag, prot := os.O_RDONLY, syscall.PROT_READ
	if writable {
		flag, prot = os.O_RDWR, syscall.PROT_READ|syscall.PROT_WRITE
	}
	f, err := os.OpenFile(path, flag, 0666)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if fi.Size() > 0 {
		return syscall.Mmap(int(f.Fd()), 0, int(fi.Size()), prot, syscall.MAP_SHARED)
	}
	return nil, nil
}

// Close releases the data. Don't read the data after running this operation
// otherwise your f*cked.
func Close(data []byte) error {
	if len(data) > 0 {
		return syscall.Munmap(data)
	}
	return nil
}
