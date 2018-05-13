package main

import (
	"os"

	"github.com/rs/zerolog"
)

func main() {
	logFile, _ := os.OpenFile("test.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)

	logger := zerolog.New(logFile).With().Logger()

	logger.Info().Str("foo", "bar").Msg("Hello World")

}
