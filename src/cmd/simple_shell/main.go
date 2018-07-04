package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		// Read the keyword input.
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		// Handle the execution of the input.
		err = execInput(input)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

var ErrNoPath = errors.New("path required")

func execInput(input string) error {
	input = strings.TrimSuffix(input, "\n")

	args := strings.Split(input, " ")

	// check for built-in commands
	switch args[0] {
	case "cd":
		if len(args) < 2 {
			return ErrNoPath
		}
		err := os.Chdir(args[1])
		if err != nil {
			return err
		}
		return nil
	case "exit":
		os.Exit(0)
	}

	cmd := exec.Command(args[0], args[1:]...) // 执行命令

	// set the correct output device.
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
