package internal

import (
	"errors"
)

var (
	ErrNotFound          = errors.New("Not Found")
	ErrDeletion          = errors.New("Type Deletion")
	ErrTableFileMagic    = errors.New("not an sstable (bad magic number)")
	ErrTableFileTooShort = errors.New("file is too short to be an sstable")
)
