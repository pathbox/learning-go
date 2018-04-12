package main

import (
	"math/rand"
	"strconv"
	"testing"
)

var NumItems int = 1000000

func BenchmarkMapStringKeys(b *testing.B) {
	m := make(map[string]string)
	k := make([]string, 0)

	for i := 0; i < NumItems; i++ {
		key := strconv.Itoa(rand.Intn(NumItems))
		m[key] = "value" + strconv.Itoa(i)
		k = append(k, key)
	}

	i := 0
	l := len(m)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if _, ok := m[k[i]]; ok {
		}

		i++
		if i >= l {
			i = 0
		}
	}
}

func BenchmarkMapIntKeys(b *testing.B) {
	m := make(map[int]string)
	k := make([]int, 0)

	for i := 0; i < NumItems; i++ {
		key := rand.Intn(NumItems)
		m[key] = "value" + strconv.Itoa(i)
		k = append(k, key)
	}

	i := 0
	l := len(m)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if _, ok := m[k[i]]; ok {
		}

		i++
		if i >= l {
			i = 0
		}
	}
}

// go test -bench=. -benchmem
// BenchmarkMapStringKeys-4        20000000               106 ns/op               0 B/op          0 allocs/op
// BenchmarkMapIntKeys-4           20000000                77.1 ns/op             0 B/op          0 allocs/op
