package main

import (
	"sync"
	"testing"
)

type Book struct {
	Title    string
	Author   string
	Pages    int
	Chapters []string
}

var pool = sync.Pool{
	New: func() interface{} {
		return &Book{}
	},
}

func BenchmarkNoPool(b *testing.B) {
	var book *Book

	for n := 0; n < b.N; n++ {
		book = &Book{
			Title:  "The Art of Computer Programming, Vol. 1",
			Author: "Donald E. Knuth",
			Pages:  672,
		}
	}

	_ = book
}

func BenchmarkPool(b *testing.B) {
	for n := 0; n < b.N; n++ {
		book := pool.Get().(*Book)
		book.Title = "The Art of Computer Programming, Vol. 1"
		book.Author = "Donald E. Knuth"
		book.Pages = 672

		pool.Put(book)
	}
}

// BenchmarkNoPool-4       30000000                45.7 ns/op            64 B/op          1 allocs/op
// BenchmarkPool-4         50000000                26.4 ns/op             0 B/op          0 allocs/op
