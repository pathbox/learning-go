package main

import (
	"strconv"
	"testing"
)

func BenchmarkIfTranfer(b *testing.B) {
	ary := make([]int, 1024)
	s := []string{}
	for n := 0; n < b.N; n++ {
		for i := 0; i < 512; i++ { // 现在只取数组前512位
			if ary[i] == 0 {
				s = append(s, "0")
			} else if ary[i] == 1 {
				s = append(s, "1")
			}
		}
	}
}

func BenchmarkParseTranfer(b *testing.B) {
	ary := make([]int, 1024)
	s := []string{}
	for n := 0; n < b.N; n++ {
		for i := 0; i < 512; i++ { // 现在只取数组前512位
			if ary[i] == 0 || ary[i] == 1 {
				s = append(s, strconv.Itoa(ary[i]))
			}
		}
	}
}

// go test -bench=. -benchmem
