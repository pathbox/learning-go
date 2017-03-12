package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Action = func(c *cli.Context) error {
		fmt.Printf("Hello %q", c.Args().Get(0))
		return nil
	}

	app.Run(os.Args)
}

// ./example2 World

// Hello World
