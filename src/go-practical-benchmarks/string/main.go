package main

import (
	"bytes"
	"strings"
	"testing"
)

var strLen int = 1000

func BenchmarkConcatString(b *testing.B) {
	var str string
	i := 0
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		str += "x"

		i++
		if i >= strLen {
			i = 0
			str = ""
		}
	}
}

func BenchmarkConcatBuffer(b *testing.B) {
	var buffer bytes.Buffer

	i := 0
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		buffer.WriteString("x")

		i++
		if i >= strLen {
			i = 0
			buffer = bytes.Buffer{}
		}
	}
}

func BenchmarkConcatBuilder(b *testing.B) {
	var builder strings.Builder

	i := 0

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		builder.WriteString("x")

		i++
		if i >= strLen {
			i = 0
			builder = strings.Builder{}
		}
	}
}

// go test -bench=. -benchmem