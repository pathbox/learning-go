package benchmark

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
)

func BenchmarkZerolog(b *testing.B) {
	log := zerolog.New(os.Stdout).With().Str("foo", "bar").Logger()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Info().
				Str("foo", "bar").
				Str("bar", "baz").
				Str("bar1", "baz").
				Str("bar2", "baz").
				Str("bar3", "baz").
				Str("bar4", "baz").
				Str("bar5", "baz").
				Str("bar6", "baz").
				Str("bar7", "baz").
				Int("n", 1).
				Msg("hello world")
		}

	})

}
