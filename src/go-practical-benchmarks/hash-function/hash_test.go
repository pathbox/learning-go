package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"math/rand"
	"testing"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/sha3"
)

func benchmarkHash(b *testing.B, h hash.Hash) {
	data := make([]byte, 1024)
	rand.Read(data)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		h.Write(data)
		h.Sum(nil)
	}
}

func BenchmarkMD5(b *testing.B) {
	benchmarkHash(b, md5.New())
}

func BenchmarkSHA1(b *testing.B) {
	benchmarkHash(b, sha1.New())
}

func BenchmarkSHA256(b *testing.B) {
	benchmarkHash(b, sha256.New())
}

func BenchmarkSHA512(b *testing.B) {
	benchmarkHash(b, sha512.New())
}

func BenchmarkSHA3256(b *testing.B) {
	benchmarkHash(b, sha3.New256())
}

func BenchmarkSHA3512(b *testing.B) {
	benchmarkHash(b, sha3.New512())
}

func BenchmarkBLAKE2b256(b *testing.B) {
	h, _ := blake2b.New256(nil)
	benchmarkHash(b, h)
}

func BenchmarkBLAKE2b512(b *testing.B) {
	h, _ := blake2b.New512(nil)
	benchmarkHash(b, h)
}

/*
go test -bench=. -benchmem
BenchmarkMD5-4                   1000000              1985 ns/op              16 B/op          1 allocs/op
BenchmarkSHA1-4                  1000000              1816 ns/op              32 B/op          1 allocs/op
BenchmarkSHA256-4                 300000              3959 ns/op              32 B/op          1 allocs/op
BenchmarkSHA512-4                 500000              2766 ns/op              64 B/op          1 allocs/op
BenchmarkSHA3256-4                300000              4773 ns/op             512 B/op          3 allocs/op
BenchmarkSHA3512-4                200000              8103 ns/op             576 B/op          3 allocs/op
BenchmarkBLAKE2b256-4            1000000              1406 ns/op              32 B/op          1 allocs/op
BenchmarkBLAKE2b512-4            1000000              1432 ns/op              64 B/op          1 allocs/op