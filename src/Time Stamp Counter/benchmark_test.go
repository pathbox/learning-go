package benchmark

import (
	"testing"
	"time"

	"github.com/templexxx/tsc"
)

func BenchmarkTscTime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tsc.UnixNano()
	}
}

func BenchmarkTime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		time.Now().UnixNano()
	}
}
