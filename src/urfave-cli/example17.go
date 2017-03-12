package main

import (
	"os"

	"github.com/urfave/cli"
)

func main() {
	cli.VersionFlag = cli.BoolFlag{
		Name:  "print-version, V",
		Usage: "print only the version",
	}

	app := cli.NewApp()
	app.Name = "partay"
	app.Version = "19.99.0"
	app.Run(os.Args)
}
