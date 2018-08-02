package ants

import (
	"errors"
	"math"
)

const (
	// DefaultPoolSize is the default capacity for a default goroutine pool
	DefaultPoolSize = math.MaxInt32

	// DefaultCleanIntervalTime is the interval time to clean up goroutines
	DefaultCleanIntervalTime = 10
)

// Init a instance pool when importing ants
var defaultPool, _ = NewPool(DefaultPoolSize)

// Submit submit a task to pool
func Submit(task f) error {
	return defaultPool.Submit(task)
}

// Running returns the number of the currently running goroutines
func Running() int {
	return defaultPool.Running()
}

// Cap returns the capacity of this default pool
func Cap() int {
	return defaultPool.Cap()
}

// Free returns the available goroutines to work
func Free() int {
	return defaultPool.Free()
}

// Release Closed the default pool
func Release() {
	defaultPool.Release()
}

// Errors for the Ants API
var (
	ErrInvalidPoolSize   = errors.New("invalid size for pool")
	ErrInvalidPoolExpiry = errors.New("invalid expiry for pool")
	ErrPoolClosed        = errors.New("this pool has been closed")
)
