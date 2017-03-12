package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

var (
	Revision = "fafafaf"
)

func main() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("version=%s revision=%s\n", c.App.Version, Revision)
	}

	app := cli.NewApp()
	app.Name = "partay"
	app.Version = "19.99.0"
	app.Run(os.Args)
}
