// You can go a long way with the Read() way of reading files, but sometimes you need more
// convenience. Something that gets used very often in Ruby are the IO functions like
//  each_line, each_char, each_codepoint etc. We can achieve something similar by using the
// Scanner type, and associated functions from the bufio package.
// The bufio.Scanner type implements functions that take a “split” function, and advance a
// pointer based on this function. For instance, the built-in bufio.ScanLines split function,
// for every iteration, advances the pointer until the next newline character. At each step,
// the type also exposes methods to obtain the byte array/string between the start and the end position.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	file, err := os.Open("filetoread.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	// Returns a boolean based on whether there's a next instance of `\n`
	// character in the IO stream. This step also advances the internal pointer
	// to the next position (after '\n') if it did find that token.
	for {
		read := scanner.Scan()

		if read { // 只是读一次的操作
			fmt.Println("read byte array: ", scanner.Bytes())
			fmt.Println("read string: ", scanner.Text())
		} else {
			break
		}
	}

}
