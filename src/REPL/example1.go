package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// cliName is the name used in the repl prompts
var cliName string = "simpleREPL"

// printPrompt displays the repl prompt at the start of each loop
func printPrompt() {
	fmt.Print(cliName, "> ")
}

// printUnkown informs the user about invalid commands
func printUnknown(text string) {
	fmt.Println(text, ": command not found")
}

// displayHelp informs the user about our hardcoded functions
func displayHelp() {
	fmt.Printf(
		"Welcome to %v! These are the available commands: \n",
		cliName,
	)
	fmt.Println(".help    - Show available commands")
	fmt.Println(".clear   - Clear the terminal screen")
	fmt.Println(".exit    - Closes your connection to ", cliName)
}

// clearScreen clears the terminal screen
func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// handleInvalidCmd attempts to recover from a bad command
func handleInvalidCmd(text string) {
	defer printUnknown(text)
}

// handleCmd parses the given commands
func handleCmd(text string) {
	handleInvalidCmd(text)
}

// cleanInput preprocesses input to the db repl
func cleanInput(text string) string {
	output := strings.TrimSpace(text)
	output = strings.ToLower(output)
	return output
}
func main() {
	// Hardcoded repl commands
	commands := map[string]interface{}{
		".help":  displayHelp,
		".clear": clearScreen,
	}
	// Begin the repl loop
	reader := bufio.NewScanner(os.Stdin)
	printPrompt()
	for reader.Scan() {
		text := cleanInput(reader.Text())
		if command, exists := commands[text]; exists {
			// Call a hardcoded function
			command.(func())()
		} else if strings.EqualFold(".exit", text) {
			// Close the program on the exit command
			return
		} else {
			// Pass the command to the parser
			handleCmd(text)
		}
		printPrompt()
	}
	// Print an additional line if we encountered an EOF character
	fmt.Println()
}
