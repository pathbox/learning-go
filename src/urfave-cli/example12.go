package main

import (
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.BoolTFlag{
			Name:  "ginger-crouton",
			Usage: "is it in the soup?",
		},
	}
	app.Action = func(ctx *cli.Context) error {
		if !ctx.Bool("ginger-crouton") {
			return cli.NewExitError("it is not in the soup", 86)
		}
		return nil
	}

	app.Run(os.Args)
}
