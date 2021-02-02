package main

import (
    "bufio"
    "fmt"
		"os"
		"strings"
)

func main() {
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        line := scanner.Text()
        if line == "exit" {
            os.Exit(0)
				}
				arr := strings.Fields(line)
        fmt.Println(arr) // Println will add back the final '\n'
    }
    if err := scanner.Err(); err != nil {
        fmt.Fprintln(os.Stderr, "reading standard input:", err)
    }
}