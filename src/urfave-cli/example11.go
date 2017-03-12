package main

import (
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Commands = []cli.Command{
		{
			Name: "noop",
		},
		{
			Name:     "add",
			Category: "template",
		},
		{
			Name:     "remove",
			Category: "template",
		},
	}

	app.Run(os.Args)
}
