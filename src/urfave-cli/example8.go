package main

import (
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "lang, l",
			Value:  "english",
			Usage:  "language for the greeting",
			EnvVar: "APP_LANG",
		},
	}

	app.Run(os.Args)
}

// Values from the Environment

// You can also have the default value set from the environment via EnvVar. e.g.
