package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
  env := os.Environ()
  procAttr := &os.ProcAttr{
    Env: env,
    Files: []*os.File{
      os.Stdin,
      os.Stdout,
      os.Stderr,
    },
  }

  pid, err := os.StartProcess("/bin/ls", []string{"ls", "-l"}, procAttr)
  if err != nil {
		fmt.Printf("Error %v starting process!\n", err)
		os.Exit(1)
	}

  fmt.Printf("The process id is %v\n", pid)

  pid, err = os.StartProcess("/bin/ps", []string{"-e", "-opid, ppid,comm"}, procAttr)
  if err != nil {
		fmt.Printf("Error %v starting process!\n", err)
		os.Exit(1)
	}
  fmt.Printf("The process id is %v\n", pid)

  cmd := exec.Command("ls", "-l")
  err = cmd.Run()
  if err != nil {
    fmt.Printf("Error %v starting process!\n", err)
    os.Exit(1)
  }
}
