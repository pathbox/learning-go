package main

import (
	"os"

	"github.com/urfave/cli"
)

func main() {
	cli.BashCompletionFlag = cli.BoolFlag{
		Name:   "compgen",
		Hidden: true,
	}

	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Commands = []cli.Command{
		{
			Name: "wat",
		},
	}
	app.Run(os.Args)
}
