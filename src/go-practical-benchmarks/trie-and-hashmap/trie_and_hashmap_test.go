package main

import (
	"fmt"
	"testing"

	trie "github.com/derekparker/trie"
)

func BenchmarkHashmapWrite(b *testing.B) {
	m := make(map[string]string)
	for n := 0; n < b.N; n++ {
		k := "/users/1/show"
		m[k] = k
	}
}

func BenchmarkTrieWrite(b *testing.B) {
	t := trie.New()
	for n := 0; n < b.N; n++ {
		k := "/users/1/show"
		t.Add(k, k)
	}
}

func BenchmarkHashmapGet(b *testing.B) {
	m := make(map[string]string)
	k := "/users/1/show/go/good"
	m[k] = k
	var r interface{}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		r = m[k]
	}
	fmt.Println(r)
}

func BenchmarkTrieGet(b *testing.B) {
	t := trie.New()
	k := "/users/1/show/go/good"
	t.Add(k, k)
	var r interface{}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		node, _ := t.Find(k)
		r = node.Meta()
	}
	fmt.Println(r)
}

// go test -bench=. -benchmem
