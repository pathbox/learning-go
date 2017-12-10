package main

import "github.com/gobwas/pool/pbytes"

func main() {
	// Reuse only slices whose capacity is 128, 256, 512 or 1024.
	pool := pbytes.New(128, 1024)

	bts := pool.GetCap(100) // Returns make([]byte, 0, 128).
	defer pool.Put(bts)

	// Work with bts.
}
