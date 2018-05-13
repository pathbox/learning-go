package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type ServerityHook struct{}

func (h ServerityHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if level != zerolog.NoLevel {
		e.Str("severity", level.String())
	}
}

func main() {
	hooked := log.Hook(ServerityHook{})
	hooked.Warn().Msg("hook test message")
}
