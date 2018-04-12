package main

import (
	"testing"
)

var numItems int = 1000000

func BenchmarkSliceAppend(b *testing.B) {
	s := make([]byte, 0)

	i := 0
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		s = append(s, 1)

		i++
		if i == numItems {
			b.StopTimer()
			i = 0
			s = make([]byte, 0)
			b.StartTimer()
		}
	}
}

func BenchmarkSliceAppendPrealloc(b *testing.B) {
	s := make([]byte, 0, numItems)

	i := 0
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		s = append(s, 1)

		i++
		if i == numItems {
			b.StopTimer()
			i = 0
			s = make([]byte, 0, numItems)
			b.StartTimer()
		}
	}
}

// BenchmarkSliceAppend-4                  1000000000               2.09 ns/op            5 B/op          0 allocs/op
// BenchmarkSliceAppendPrealloc-4          2000000000               1.15 ns/op            0 B/op          0 allocs/
