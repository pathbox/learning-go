package main

import (
	"os"

	"github.com/urfave/cli"
)

func main() {
	cli.HelpFlag = cli.BoolFlag{
		Name:   "halp, haaaaalp",
		Usage:  "HALP",
		EnvVar: "SHOW_HALP,HALPPLZ",
	}

	cli.NewApp().Run(os.Args)
}
