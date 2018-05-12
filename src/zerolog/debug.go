package main

import (
	log "github.com/rs/zerolog/log"

	"github.com/rs/zerolog"
)

func main() {
	// UNIX Time is faster and smaller than most timestamps
	// If you set zerolog.TimeFieldFormat to an empty string,
	// logs will write with UNIX time
	zerolog.TimeFieldFormat = ""
	log.Debug().
		Str("Scale", "833 cents").     // custom key and string value
		Float64("Interval", 833.09).   // custom key and float64 value
		Msg("Fibonacci is everywhere") // the message key and value
}
