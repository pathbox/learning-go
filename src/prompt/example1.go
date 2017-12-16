package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/manifoldco/promptui"
)

func main() {
	validate := func(input string) error { // input is from terminal
		_, err := strconv.ParseFloat(input, 64)
		if err != nil {
			return errors.New("Invalid number")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Number",
		Validate: validate, // validate is a func
	}

	result, err := prompt.Run() // How to get the result

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fmt.Printf("You choose %q\n", result)
}
