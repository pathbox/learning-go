package main

import (
	"bufio"
	"io"
	"os"
	"testing"
)

func BenchmarkWriteFile(b *testing.B) {
	for n := 0; n < b.N; n++ {
		f, err := os.Create("/tmp/test.txt")
		if err != nil {
			panic(err)
		}

		for i := 0; i < 100000; i++ {
			f.WriteString("Some text!\n")
		}

		f.Close()
	}
}

func BenchmarkWriteFileBuffered(b *testing.B) {
	for n := 0; n < b.N; n++ {
		f, err := os.Create("/tmp/test.txt")
		if err != nil {
			panic(err)
		}

		w := bufio.NewWriter(f)

		for i := 0; i < 100000; i++ {
			w.WriteString("some text!\n")
		}

		w.Flush()
		f.Close()
	}
}

func BenchmarkReadFile(b *testing.B) {
	for n := 0; n < b.N; n++ {
		f, err := os.Open("/tmp/test.txt")
		if err != nil {
			panic(err)
		}

		b := make([]byte, 10)

		_, err = f.Read(b)
		for err == nil {
			_, err = f.Read(b)
		}
		if err != io.EOF {
			panic(err)
		}

		f.Close()
	}
}

func BenchmarkReadFileBuffered(b *testing.B) {
	for n := 0; n < b.N; n++ {
		f, err := os.Open("/tmp/test.txt")
		if err != nil {
			panic(err)
		}

		r := bufio.NewReader(f)

		_, err = r.ReadString('\n')
		for err == nil {
			_, err = r.ReadString('\n')
		}
		if err != io.EOF {
			panic(err)
		}

		f.Close()
	}
}

/* go test -bench=. -benchmem
BenchmarkWriteFile-4                   2         697476746 ns/op             136 B/op          3 allocs/op
BenchmarkWriteFileBuffered-4         300           4991695 ns/op            4200 B/op          4 allocs/op
BenchmarkReadFile-4                   10         184628713 ns/op             104 B/op          3 allocs/op
BenchmarkReadFileBuffered-4          200          10109582 ns/op         3204201 B/op     200004 allocs/op

使用bufio 是空间换时间，会消耗更多的内存。 通过benchmark，bufio.NewWriter(f) 起到的优化效果更佳，read 方面消耗的内存过多
使用中要注意这一点
*/
