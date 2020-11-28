package netutil

import "syscall"

func maxOpenFiles() int {
	var rlim syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rlim); err != nil {
		return defaultMaxOpenFiles
	}
	return int(rlim.Cur)
}