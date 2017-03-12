package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "boom"
	app.Usage = "make an explosive entrance"
	app.Action = func(c *cli.Context) error {
		fmt.Println("boom! I say!")
		return nil
	}

	app.Run(os.Args)
}

// Install our command to the $GOPATH/bin directory:

// $ go install
// Finally run our new command:

// $ greet
// Hello friend!
// cli also generates neat help text:

// $ greet help
// NAME:
//     greet - fight the loneliness!

// USAGE:
//     greet [global options] command [command options] [arguments...]

// VERSION:
//     0.0.0

// COMMANDS:
//     help, h  Shows a list of commands or help for one command

// GLOBAL OPTIONS
//     --version Shows version information
