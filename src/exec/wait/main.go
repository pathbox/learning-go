package main

import (
	"fmt"
	"os/exec"
	"time"
)

func main() {
  cmd := exec.Command("sleep", "5")
  if err := cmd.Start(); err != nil {
    panic(err)
  }

  done := make(chan error, 1)
  go func(){
    done <- cmd.Wait()
  }()
  select{
  case <- time.After(3 * time.Second):
    cmd.Process.Kill()
    fmt.Println("timeout")
  case <- done:
    fmt.Println("done")
  }
}
