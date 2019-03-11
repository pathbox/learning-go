package mergesort

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const size = 1000000

func TestMergesort(t *testing.T) {
	s := []int{5, 8, 9, 5, 0, 10, 1, 6}
	mergesort(s)
	assert.Equal(t, []int{0, 1, 5, 5, 6, 8, 9, 10}, s)
}

func TestParallelMergesort1(t *testing.T) {
	s := []int{5, 8, 9, 5, 0, 10, 1, 6}
	parallelMergesort1(s)
	assert.Equal(t, []int{0, 1, 5, 5, 6, 8, 9, 10}, s)
}

func TestParallelMergesort2(t *testing.T) {
	s := []int{5, 8, 9, 5, 0, 10, 1, 6}
	parallelMergesort2(s)
	assert.Equal(t, []int{0, 1, 5, 5, 6, 8, 9, 10}, s)
}

func TestParallelMergesort3(t *testing.T) {
	s := []int{5, 8, 9, 5, 0, 10, 1, 6}
	parallelMergesort3(s)
	assert.Equal(t, []int{0, 1, 5, 5, 6, 8, 9, 10}, s)
}

func BenchmarkMergesort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := random(size)
		b.StartTimer()
		mergesort(s)
		b.StopTimer()
	}
}

func BenchmarkParallelMergesort1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := random(size)
		b.StartTimer()
		parallelMergesort1(s)
		b.StopTimer()
	}
}

func BenchmarkParallelMergesort2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := random(size)
		b.StartTimer()
		parallelMergesort2(s)
		b.StopTimer()
	}
}

func BenchmarkParallelMergesort3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := random(size)
		b.StartTimer()
		parallelMergesort3(s)
		b.StopTimer()
	}
}
