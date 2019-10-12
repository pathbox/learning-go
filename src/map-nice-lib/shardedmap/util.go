package shardedmap

import (
	"runtime"
	"unsafe"
)

//nolint:gochecknoglobals
var defaultShards = runtime.NumCPU() * 16 // github.com/tidwall/shardmap recommendation

// Adapted from https://github.com/dgraph-io/ristretto/blob/master/z/rtutil.go
//
// MIT License
// Copyright (c) 2019 Ewan Chou
//
// Not copying the whole thing as this repo itself is under MIT License. If
// that's considered a violation, just message me.

//go:noescape
//go:linkname rtmemhash runtime.memhash
func rtmemhash(p unsafe.Pointer, h, s uintptr) uintptr

type stringStruct struct {
	str unsafe.Pointer
	len int
}

// memHash is the hash function used by go map, it utilizes available hardware instructions(behaves
// as aeshash if aes instruction is available).
// NOTE: The hash seed changes for every process. So, this cannot be used as a persistent hash.
func memHash(data []byte) uint64 {
	ss := (*stringStruct)(unsafe.Pointer(&data))
	return uint64(rtmemhash(ss.str, 0, uintptr(ss.len)))
}

// memHashString is the hash function used by go map, it utilizes available hardware instructions
// (behaves as aeshash if aes instruction is available).
// NOTE: The hash seed changes for every process. So, this cannot be used as a persistent hash.
func memHashString(str string) uint64 {
	ss := (*stringStruct)(unsafe.Pointer(&str))
	return uint64(rtmemhash(ss.str, 0, uintptr(ss.len)))
}