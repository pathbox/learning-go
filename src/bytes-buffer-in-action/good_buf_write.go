package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func process(data []string) string {
	var buf bytes.Buffer

	for _, item := range data {
		buf.WriteString(item)
		buf.WriteString("\n")
	}

	return buf.String()
}

func BenchmarkOld(b *testing.B) {
	data := make([]string, 100000)

	for i := 0; i < len(data); i++ {
		data[i] = fmt.Sprintf("Item %d", i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := process(data)
		_ = result
	}
}

// 复用buf
func newProcess(data []string, buf *bytes.Buffer) string {
	buf.Reset()

	for _, item := range data {
		buf.WriteString(item)
		buf.WriteString("\n")
	}

	return buf.String()
}

func BenchmarkNew(b *testing.B) {
	data := make([]string, 100000)
	for i := 0; i < len(data); i++ {
		data[i] = fmt.Sprintf("Item %d", i)
	}

	buf := new(bytes.Buffer)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result := newProcess(data, buf)
		_ = result
	}
}

func strBuilderProcess(data []string, buf *strings.Builder) string {
	buf.Reset()

	for _, item := range data {
		buf.WriteString(item)
		buf.WriteString("\n")
	}

	return buf.String()
}

func BenchmarkStr(b *testing.B) {
	data := make([]string, 100000)
	for i := 0; i < len(data); i++ {
		data[i] = fmt.Sprintf("Item %d", i)
	}

	buf := new(strings.Builder)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := strBuilderProcess(data, buf)
		_ = result
	}
}

func BenchmarkJoin(b *testing.B) {
	data := make([]string, 100000)
	for i := 0; i < len(data); i++ {
		data[i] = fmt.Sprintf("Item %d", i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := strings.Join(data, "\n")
		_ = result
	}
}

//  join 和 复用 bytes.Buffer 性能最好
