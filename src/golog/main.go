package main

import (
    "bufio"
    "log"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "time"
)

func GetCurrentPath() string  {
  file, _ := exec.LookPath(os.Args[0])
  path, _ := filepath.Abs(file)
  splitstring := strings.Split(path, "\\")
  size := len(splitstring)
  splitstring = strings.Split(path, splitstring[size-1])
  ret := strings.Replace(splitstring[0], "\\", "/", size-1)
  return ret
}

func main() {
  t := time.Now()
  filepath := "./log_" + strings.Replace(t.String()[:19], ":", "_", 3) + ".log"
  file, err := os.OpenFile(filepath, os.O_CREATE, 0666)
  if err != nil {
    log.Fatal("create log file failed")
  }
  defer file.Close()
  wFile := bufio.NewWriter(file)
  wFile.WriteString(GetCurrentPath())
  wFile.Flush()
}
