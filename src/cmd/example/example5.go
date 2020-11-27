package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	os.Setenv("NAME", "wangbm")
	cmd := exec.Command("echo", os.ExpandEnv("$NAME"))
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	fmt.Printf("%s", out)
}
