package iptables

import (
	"os"
	"sync"
	"syscall"
)

const (
	// In earlier versions of iptables, the xtables lock was implemented
	// via a Unix socket, but now flock is used via this lockfile:
	// http://git.netfilter.org/iptables/commit/?id=aa562a660d1555b13cffbac1e744033e91f82707
	// Note the LSB-conforming "/run" directory does not exist on old
	// distributions, so assume "/var" is symlinked
	xtablesLockFilePath = "/var/run/xtables.lock"

	defaultFilePerm = 0600
)

type Unlocker interface {
	Unlock() error
}

type nopUnlocker struct{}

func (_ nopUnlocker) Unlock() error { return nil }

type fileLock struct {
	mu sync.Mutex
	fd int
}

// tryLock takes an exclusive lock on the xtables lock file without blocking.
// This is best-effort only: if the exclusive lock would block (i.e. because
// another process already holds it), no error is returned. Otherwise, any
// error encountered during the locking operation is returned.
// The returned Unlocker should be used to release the lock when the caller is
// done invoking iptables commands.
func (l *fileLock) tryLock() (Unlocker, error) {
	l.mu.Lock()
	err := syscall.Flock(l.fd, syscall.LOCK_EX|syscall.LOCK_NB)
	switch err {
	case syscall.EWOULDBLOCK:
		l.mu.Unlock()
		return nopUnlocker{}, nil
	case nil:
		return l, nil
	default:
		l.mu.Unlock()
		return nil, err
	}
}

func (l *fileLock) Unlock() error {
	defer l.mu.Unlock()
	return syscall.Close(l.fd)
}

func newXtablesFileLock() (*fileLock, error) {
	fd, err := syscall.Open(xtablesLockFilePath, os.O_CREATE, defaultFilePerm)
	if err != nil {
		return nil, err
	}
	return &fileLock{fd: fd}, nil
}