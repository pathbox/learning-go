package main

import (
	log "github.com/rs/zerolog/log"

	"github.com/rs/zerolog"
)

func main() {
	zerolog.TimestampFieldName = "t"
	zerolog.LevelFieldName = "l"
	zerolog.MessageFieldName = "m"
	// zerolog.TimeFieldFormat = time.RFC3339 // 默认打印的为时间，定义空串为时间戳，速度最快
	// zerolog.TimeFieldFormat = ""

	log.Info().Msg("hello world")
}
