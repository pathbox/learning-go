package main

import (
	"github.com/rs/zerolog"
	log "github.com/rs/zerolog/log"
)

func main() {
	log.Info().
		Str("foo", "bar").
		Dict("dict", zerolog.Dict().
			Str("bar", "baz").
			Int("n", 1),
		).Msg("hello world")
}
