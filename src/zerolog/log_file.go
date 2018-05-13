package main

import (
	"os"

	"github.com/rs/zerolog"
)

func main() {
	log := zerolog.New(os.Stdout).With().
		Str("foo", "bar").
		Logger()

	stdlog.SetFlags(0)
	stdlog.SetOutput(log)

	stdlog.Print("hello world")
}
