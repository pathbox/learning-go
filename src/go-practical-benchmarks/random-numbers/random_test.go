package main

import (
	crand "crypto/rand"
	"math/big"
	"math/rand"
	"testing"
)

func BenchmarkMathRand(b *testing.B) {
	for n := 0; n < b.N; n++{
		rand.Int63()
	}
}

func BenchmarkCryptoRand(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, err := crand.Int(crand.Reader, big.NewInt(27))
		if err != nil {
			panic(err)
		}
	}
}

// go test -bench=. -benchmem

// goos: darwin
// goarch: amd64
// BenchmarkMathRand-4             50000000                27.7 ns/op             0 B/op          0 allocs/op
// BenchmarkCryptoRand-4            3000000               426 ns/op             161 B/op          5 allocs/op
// PASS