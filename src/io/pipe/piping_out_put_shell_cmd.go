package main

import (
	"io"
	"log"
	"os"
	"os/exec"
)

func main() {
	pr, pw := io.Pipe()
	defer pw.Close()

	cmd := exec.Command("cat", "server.go")
	cmd.Stdout = pw

	go func() {
		defer pr.Close()

		if _, err := io.Copy(os.Stdout, pr); err != nil {
			log.Fatal(err)
		}
	}()

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
