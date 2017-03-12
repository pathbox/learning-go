package main

import (
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "lang, l",
			Value: "english",
			Usage: "language for the greeting",
		},
	}

	app.Run(os.Args)
}

// That flag can then be set with --lang spanish or -l spanish. Note that giving two different forms of the same flag in the same command invocation is an error.
