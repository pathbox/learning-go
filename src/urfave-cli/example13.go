package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {
	tasks := []string{"cook", "clean", "laundry", "eat", "sleep", "code"}

	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Commands = []cli.Command{
		{
			Name:    "complete",
			Aliases: []string{"c"},
			Usage:   "complete a task on the list",
			Action: func(c *cli.Context) error {
				fmt.Println("completed task: ", c.Args().First())
				return nil
			},
			BashComplete: func(c *cli.Context) {
				// This will complete if no args are passed
				if c.NArg() > 0 {
					return
				}
				for _, t := range tasks {
					fmt.Println(t)
				}
			},
		},
	}

	app.Run(os.Args)
}
