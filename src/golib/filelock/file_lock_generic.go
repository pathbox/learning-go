package filelock

import (
	"fmt"
	"io"
	"runtime"
)

func Lock(name string) (io.Closer, error) {
	return nil, fmt.Errorf("leveldb/db: file locking is not implemented on %s/%s", runtime.GOOS, runtime.GOARCH)
}
