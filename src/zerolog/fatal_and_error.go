package main

import (
	"errors"

	"github.com/rs/zerolog"
	log "github.com/rs/zerolog/log"
)

func main() {
	err := errors.New("A repo man spends his life getting into tense situations")

	service := "myservice"

	zerolog.TimeFieldFormat = ""
	log.Error().
		Str("error", "err msg").
		Msg("error message")
	log.Fatal(). // log之后 exit 1  离开进程
			Err(err).
			Str("service", service).
			Msgf("Cannot start %s", service)
}
